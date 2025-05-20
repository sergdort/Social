package postsapp

import (
	"github.com/sergdort/Social/app/shared/mid"
	"github.com/sergdort/Social/business/domain"
	"github.com/sergdort/Social/foundation/web"
	"net/http"
)

type Config struct {
	Auth      *domain.AuthUseCase
	PostsRepo domain.PostsRepository
}

func Routes(app *web.App, config Config) {
	const version = "v1"

	api := postsApp{repo: config.PostsRepo}
	auth := mid.Bearer(config.Auth)
	postContext := api.postsContextMiddleware()

	app.HandlerFunc(http.MethodPost, version, "/posts", api.createPostsHandler, auth)
	app.HandlerFunc(http.MethodGet, version, "/posts/{postId}", api.getPostHandler, auth, postContext)
}
