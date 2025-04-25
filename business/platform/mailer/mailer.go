package mailer

type Mailer interface {
	Send(templateFile, username, email string, data any) error
}
