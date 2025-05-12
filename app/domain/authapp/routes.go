package authapp

import (
	"github.com/sergdort/Social/business/domain"
	"github.com/sergdort/Social/foundation/web"
	"net/http"
)

type Config struct {
	UseCase *domain.AuthUseCase
}

func Routes(app *web.App, config Config) {
	const version = "v1"

	api := authApp{useCase: config.UseCase}

	app.HandlerFunc(http.MethodPost, version, "/authentication/user", api.registerUserHandler)
	app.HandlerFunc(http.MethodPost, version, "/authentication/token", api.createTokenHandler)
}
