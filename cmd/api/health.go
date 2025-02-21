package main

import (
	"net/http"
)

func (app *application) healthHandler(writer http.ResponseWriter, request *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}
	if err := app.jsonResponse(writer, http.StatusOK, data); err != nil {
		app.internalServerError(writer, request, err)
	}
}
