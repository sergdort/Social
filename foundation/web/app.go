package web

import (
	"context"
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"net/http"
)

// Encoder defines behavior that can encode a data model and provide
// the content type for that encoding.
type Encoder interface {
	Encode() (data []byte, contentType string, err error)
}

// HandlerFunc represents a function that handles a http request within our own
// little mini framework.
type HandlerFunc func(ctx context.Context, r *http.Request) Encoder

// MidFunc is a handler function designed to run code before and/or after
// another Handler. It is designed to remove boilerplate or other concerns not
// direct to any given app Handler.
type MidFunc func(handler HandlerFunc) HandlerFunc

// Logger represents a function that will be called to add information
// to the logs.
type Logger func(ctx context.Context, msg string, args ...any)

type App struct {
	mux     *http.ServeMux
	mw      []MidFunc
	otmux   http.Handler
	origins []string
}

func NewApp(docsSetup func(mux *http.ServeMux), mw ...MidFunc) *App {
	mux := http.NewServeMux()
	docsSetup(mux)
	app := &App{
		mux:   mux,
		mw:    mw,
		otmux: otelhttp.NewHandler(mux, "request"),
	}
	return app
}

// ServeHTTP implements the http.Handler interface. It's the entry point for
// all http traffic and allows the opentelemetry mux to run first to handle
// tracing. The opentelemetry mux then calls the application mux to handle
// application traffic. This was set up in the NewApp function.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.otmux.ServeHTTP(w, r)
}

func (a *App) EnableCORS(origins []string) {
	a.origins = origins

	handler := func(ctx context.Context, r *http.Request) Encoder {
		return nil
	}
	handler = wrapMiddleware([]MidFunc{a.corsHandler}, handler)

	a.HandlerFuncNoMid(http.MethodOptions, "", "/", handler)
}

func (a *App) corsHandler(webHandler HandlerFunc) HandlerFunc {
	h := func(ctx context.Context, r *http.Request) Encoder {
		w := GetWriter(ctx)

		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Origin
		//
		// Limiting the possible Access-Control-Allow-Origin values to a set of
		// allowed origins requires code on the server side to check the value of
		// the Origin request header, compare that to a list of allowed origins, and
		// then if the Origin value is in the list, set the
		// Access-Control-Allow-Origin value to the same value as the Origin.

		reqOrigin := r.Header.Get("Origin")
		for _, origin := range a.origins {
			if origin == "*" || origin == reqOrigin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "POST, PATCH, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")

		return webHandler(ctx, r)
	}

	return h
}

// HandlerFunc sets a handler function for a given HTTP method and path pair
// to the application server mux.
func (a *App) HandlerFunc(method string, group string, path string, handlerFunc HandlerFunc, mw ...MidFunc) {
	handlerFunc = wrapMiddleware(mw, handlerFunc)
	handlerFunc = wrapMiddleware(a.mw, handlerFunc)

	h := func(w http.ResponseWriter, r *http.Request) {
		//ctx := setTracer(r.Context(), app.tracer)
		ctx := setWriter(r.Context(), w)

		otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(w.Header()))

		resp := handlerFunc(ctx, r)

		if err := Respond(ctx, w, resp); err != nil {
			//app.log(ctx, "web-respond", "ERROR", err)
			return
		}
	}
	finalPath := path
	if group != "" {
		finalPath = "/" + group + path
	}
	finalPath = fmt.Sprintf("%s %s", method, finalPath)

	a.mux.HandleFunc(finalPath, h)
}

// HandlerFuncNoMid sets a handler function for a given HTTP method and path
// pair to the application server mux. Does not include the application
// middleware or OTEL tracing.
func (a *App) HandlerFuncNoMid(method string, group string, path string, handlerFunc HandlerFunc) {
	h := func(w http.ResponseWriter, r *http.Request) {
		ctx := setWriter(r.Context(), w)

		resp := handlerFunc(ctx, r)

		if err := Respond(ctx, w, resp); err != nil {
			//a.log(ctx, "web-respond", "ERROR", err)
			return
		}
	}

	finalPath := path
	if group != "" {
		finalPath = "/" + group + path
	}
	finalPath = fmt.Sprintf("%s %s", method, finalPath)

	a.mux.HandleFunc(finalPath, h)
}

func (a *App) AddHocHandler(method string, group string, path string, handlerFunc http.HandlerFunc) {
	finalPath := path
	if group != "" {
		finalPath = "/" + group + path
	}
	finalPath = fmt.Sprintf("%s %s", method, finalPath)

	a.mux.HandleFunc(finalPath, handlerFunc)
}
