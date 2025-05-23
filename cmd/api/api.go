package main

import (
	"context"
	"github.com/sergdort/Social/app/domain/authapp"
	"github.com/sergdort/Social/app/domain/feedapp"
	"github.com/sergdort/Social/app/domain/postsapp"
	"github.com/sergdort/Social/app/domain/usersapp"
	"github.com/sergdort/Social/app/shared/mid"
	"github.com/sergdort/Social/business/domain"
	"github.com/sergdort/Social/business/platform/mailer"
	s "github.com/sergdort/Social/business/platform/store"
	"github.com/sergdort/Social/business/platform/store/cache"
	"github.com/sergdort/Social/docs" // This is required to generate Swagger docs
	"github.com/sergdort/Social/foundation/logger"
	"github.com/sergdort/Social/foundation/otel"
	"github.com/sergdort/Social/foundation/web"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"time"
)

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type application struct {
	config  config
	store   s.Storage
	logger  *logger.Logger
	mailer  mailer.Mailer
	cache   cache.Storage
	useCase useCases
}

type useCases struct {
	Users *domain.UsersUseCase
	Auth  *domain.AuthUseCase
	Feed  domain.FeedRepository
	Posts domain.PostsRepository
}

type redisConfig struct {
	addr    string
	pw      string
	db      int
	enabled bool
}

type config struct {
	address         string
	debugHost       string
	shutDownTimeout time.Duration
	db              dbConfig
	env             string
	apiURL          string
	mail            mailConfig
	frontEndURL     string
	auth            authConfig
	redisCfg        redisConfig
	serviceName     string
}

type mailConfig struct {
	exp            time.Duration
	fromEmail      string
	sendGridConfig sendGridConfig
}

type authConfig struct {
	basic basicAuthConfig
	jwt   jwtAuthConfig
}
type basicAuthConfig struct {
	username string
	password string
}

type jwtAuthConfig struct {
	secret    string
	exp       time.Duration
	tokenHost string
}

type sendGridConfig struct {
	apiKey string
}

func (app *application) mount(ctx context.Context, log *logger.Logger) http.Handler {
	traceProvider, teardown, err := otel.InitTracing(log, otel.Config{
		ServiceName: app.config.frontEndURL,
		Host:        app.config.debugHost,
		ExcludedRoutes: map[string]struct{}{
			"/v1/liveness":  {},
			"/v1/readiness": {},
		},
		Probability: 0.05,
	})
	if err != nil {
		log.Error(ctx, "starting tracing: %s", err)
		return nil
	}
	tracer := traceProvider.Tracer(app.config.serviceName)
	webApp := web.NewApp(
		func(mux *http.ServeMux) {
			mux.HandleFunc("GET /v1/swagger/", httpSwagger.Handler())
		},
		mid.Otel(tracer),
		mid.Logger(log),
		mid.Errors(log),
		mid.Metrics(),
		mid.Panics(),
	)

	webApp.EnableCORS([]string{})
	//
	//router.Route("/v1", func(r chi.Router) {
	//	r.With(app.BasicAuthMiddleware).Get("/health", app.healthHandler)
	//	r.Get("/swagger/*", httpSwagger.Handler())
	//
	//	r.Route("/posts", func(r chi.Router) {
	//		r.Use(app.AuthTokenMiddleware)
	//		r.Post("/", app.createPostsHandler)
	//
	//		r.Route("/{postID}", func(r chi.Router) {
	//			r.Use(app.postsContextMiddleware)
	//			r.Get("/", app.getPostHandler)
	//			r.Delete("/", app.checkPostOwnershipMiddleware(domain.RoleTypeAdmin, app.deletePostHandler))
	//			r.Patch("/", app.checkPostOwnershipMiddleware(domain.RoleTypeModerator, app.patchPostsHandler))
	//		})
	//	})
	//
	//	r.Route("/users", func(r chi.Router) {
	//		r.Put("/activate/{token}", app.activateUserHandler)
	//		r.Route("/{userID}", func(r chi.Router) {
	//			r.Use(app.AuthTokenMiddleware)
	//			r.Use(app.userContextMiddleware)
	//			//r.Get("/", app.getUserHandler)
	//			r.Put("/follow", app.followUserHandler)
	//			r.Put("/unfollow", app.unfollowUserHandler)
	//		})
	//
	//		r.Group(func(r chi.Router) {
	//			r.Get("/feed", app.getUserFeedHandler)
	//		})
	//	})
	//
	//	r.Route("/authentication", func(r chi.Router) {
	//		r.Post("/user", app.registerUserHandler)
	//		r.Post("/token", app.createTokenHandler)
	//	})
	//})

	authapp.Routes(webApp, authapp.Config{UseCase: app.useCase.Auth})
	usersapp.Routes(webApp, usersapp.Config{Auth: app.useCase.Auth, UseCase: app.useCase.Users})
	postsapp.Routes(webApp, postsapp.Config{Auth: app.useCase.Auth, PostsRepo: app.useCase.Posts})
	feedapp.Routes(webApp, feedapp.Config{Auth: app.useCase.Auth, FeedUseCase: app.useCase.Feed})
	defer teardown(ctx)

	return webApp
}

func (app *application) makeServer(mux http.Handler) *http.Server {
	// Docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.BasePath = "/v1"

	server := &http.Server{
		Addr:         app.config.address,
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}
	return server
}
