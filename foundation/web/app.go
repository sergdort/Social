package web

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"net/http"
)

type App struct {
	MUX *chi.Mux
}

func NewApp(mux *chi.Mux) *App {
	app := &App{
		mux,
	}
	return app
}

func (app *App) EnableCORS() {
	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	app.MUX.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
}

func (app *App) UseMiddleware(middlewares ...func(http.Handler) http.Handler) {
	app.MUX.Use(middlewares...)
}
