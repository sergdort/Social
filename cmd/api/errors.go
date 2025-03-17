package main

import (
	"net/http"
	"runtime"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	pc, _, _, ok := runtime.Caller(1) // 1 means "who called this function"
	caller := "unknown"
	if ok {
		caller = runtime.FuncForPC(pc).Name() // Get function name
	}
	app.logger.Errorw("Internal error", "caller", caller, "method", r.Method, "path", r.URL.Path, "error", err.Error())
	_ = writeJSONError(w, http.StatusInternalServerError, "Internal Server Error")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	pc, _, _, ok := runtime.Caller(1) // 1 means "who called this function"
	caller := "unknown"
	if ok {
		caller = runtime.FuncForPC(pc).Name() // Get function name
	}
	app.logger.Warnw("Bad request", "caller", caller, "method", r.Method, "path", r.URL.Path, "error", err.Error())
	_ = writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	pc, _, _, ok := runtime.Caller(1) // 1 means "who called this function"
	caller := "unknown"
	if ok {
		caller = runtime.FuncForPC(pc).Name() // Get function name
	}
	app.logger.Warnw("Not found", "caller", caller, "method", r.Method, "path", r.URL.Path, "error", err.Error())
	_ = writeJSONError(w, http.StatusNotFound, "Not found")
}

func (app *application) unauthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	pc, _, _, ok := runtime.Caller(1)
	caller := "unknown"
	if ok {
		caller = runtime.FuncForPC(pc).Name() // Get function name
	}
	app.logger.Warnw("Unauthorized Basic", "caller", caller, "method", r.Method, "path", r.URL.Path, "error", err.Error())
	w.Header().Set("WWW-Authenticate", `Basic realm="api", charset="UTF-8"`)
	_ = writeJSONError(w, http.StatusUnauthorized, err.Error())
}

func (app *application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	pc, _, _, ok := runtime.Caller(1)
	caller := "unknown"
	if ok {
		caller = runtime.FuncForPC(pc).Name() // Get function name
	}
	app.logger.Warnw("Unauthorized", "caller", caller, "method", r.Method, "path", r.URL.Path, "error", err.Error())
	_ = writeJSONError(w, http.StatusUnauthorized, err.Error())
}
