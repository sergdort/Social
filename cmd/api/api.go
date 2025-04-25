package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/sergdort/Social/business/domain"
	"github.com/sergdort/Social/docs" // This is required to generate Swagger docs
	"github.com/sergdort/Social/internal/auth"
	"github.com/sergdort/Social/internal/mailer"
	s "github.com/sergdort/Social/internal/store"
	"github.com/sergdort/Social/internal/store/cache"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
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
	config        config
	store         s.Storage
	logger        *zap.SugaredLogger
	mailer        mailer.Mailer
	authenticator auth.Authenticator
	cache         cache.Storage
	useCase       useCases
}

type useCases struct {
	Users *domain.UsersUseCase
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

func (app *application) mount() http.Handler {
	var router = chi.NewRouter()

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RealIP)
	router.Use(middleware.RequestID)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(60 * time.Second))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	router.Route("/v1", func(r chi.Router) {
		r.With(app.BasicAuthMiddleware).Get("/health", app.healthHandler)
		r.Get("/swagger/*", httpSwagger.Handler())

		r.Route("/posts", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Post("/", app.createPostsHandler)

			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)
				r.Get("/", app.getPostHandler)
				r.Delete("/", app.checkPostOwnershipMiddleware(domain.RoleTypeAdmin, app.deletePostHandler))
				r.Patch("/", app.checkPostOwnershipMiddleware(domain.RoleTypeModerator, app.patchPostsHandler))
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.activateUserHandler)
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Use(app.userContextMiddleware)
				r.Get("/", app.getUserHandler)
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})

			r.Group(func(r chi.Router) {
				r.Get("/feed", app.getUserFeedHandler)
			})
		})

		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler)
			r.Post("/token", app.createTokenHandler)
		})
	})

	return router
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

	app.logger.Infow(
		"Server has started",
		"addr", app.config.address,
		"env", app.config.env,
	)

	return server
}
