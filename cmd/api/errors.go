package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Internal server error method: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())
	_ = writeJSONError(w, http.StatusInternalServerError, "Internal Server Error")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Bad Request error method: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())
	_ = writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Not found method: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())
	_ = writeJSONError(w, http.StatusNotFound, "Not found")
}
