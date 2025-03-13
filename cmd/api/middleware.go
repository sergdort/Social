package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func (app *application) BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("missing authorization header"))
			return
		}
		// "Basic tokenstringhere"
		base64Token := authHeader[6:]

		decodedToken, err := base64.StdEncoding.DecodeString(base64Token)
		if err != nil {
			app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("invalid base64 token"))
			return
		}

		username := app.config.auth.basic.username
		pass := app.config.auth.basic.password

		creds := strings.SplitN(string(decodedToken), ":", 2)
		if len(creds) != 2 || creds[0] != username || creds[1] != pass {
			app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("invalid token"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
