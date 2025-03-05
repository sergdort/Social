package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"html/template"
	"log"
	"time"
)

const (
	FromName            = "Social"
	MaxRetries          = 3
	UserWelcomeTemplate = "user_invitation.tmpl"
)

//go:embed "templates"
var FS embed.FS

type SendgridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendgridMailer(fromEmail, apiKey string) Mailer {
	return &SendgridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    sendgrid.NewSendClient(apiKey),
	}
}

func (m *SendgridMailer) Send(templateFile, username, email string, data any) error {
	from := mail.NewEmail(FromName, m.fromEmail)
	to := mail.NewEmail(username, email)
	// Template parsing and building

	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return err
	}
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	body := new(bytes.Buffer)

	err = tmpl.ExecuteTemplate(body, "body", data)

	message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())

	for i := 0; i < MaxRetries; i++ {
		_, err := m.client.Send(message)
		if err != nil {
			log.Printf("Error sending email: %v, attempt %d", email, i+1)
			log.Printf("Error: %v", err.Error())

			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		return nil
	}

	return fmt.Errorf("error sending email")
}
