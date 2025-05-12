package usersapp

import (
	"github.com/sergdort/Social/business/domain"
	"github.com/sergdort/Social/foundation/web"
	"net/http"
)

type Config struct {
	UseCase *domain.UsersUseCase
}

func Routes(app *web.App, config Config) {
	const version = "v1"

	api := newApp(config.UseCase)
	userContext := api.userContextMiddleware(config.UseCase)
	app.HandlerFunc(http.MethodGet, version, "/users/{userID}", api.getUserHandler, userContext)
}
