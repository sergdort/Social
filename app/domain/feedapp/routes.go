package feedapp

import (
	"github.com/sergdort/Social/app/shared/mid"
	"github.com/sergdort/Social/business/domain"
	"github.com/sergdort/Social/foundation/web"
	"net/http"
)

type Config struct {
	Auth        *domain.AuthUseCase
	FeedUseCase domain.FeedUseCase
}

func Routes(app *web.App, config Config) {
	const version = "v1"

	api := feedApp{feedUseCase: config.FeedUseCase}
	auth := mid.Bearer(config.Auth)

	app.HandlerFunc(http.MethodGet, version, "/user/feed", api.getFeedHandler, auth)
}
