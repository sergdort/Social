package usersapp

import (
	"github.com/sergdort/Social/app/shared/mid"
	"github.com/sergdort/Social/business/domain"
	"github.com/sergdort/Social/foundation/web"
	"net/http"
)

type Config struct {
	Auth    *domain.AuthUseCase
	UseCase *domain.UsersUseCase
}

func Routes(app *web.App, config Config) {
	const version = "v1"

	api := newApp(config.UseCase)
	auth := mid.Bearer(config.Auth)
	userContext := api.userContextMiddleware(config.UseCase)

	app.HandlerFunc(http.MethodGet, version, "/users/{userID}", api.getUserHandler, auth, userContext)
	app.HandlerFunc(http.MethodPut, version, "/users/{userID}/follow", api.followUserHandler, auth, userContext)
}
