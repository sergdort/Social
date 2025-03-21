package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sergdort/Social/internal/store"
	"net/http"
	"strconv"
	"strings"
)

type authUserKey string

const authUserCtx userKey = "authUser"

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

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("missing authorization header"))
			return
		}
		// "Bearer <token>"
		token := authHeader[7:]

		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}
		claims := jwtToken.Claims.(jwt.MapClaims)
		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}
		ctx := r.Context()
		user, err := app.store.Users.GetByID(ctx, userID)
		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}
		ctx = context.WithValue(ctx, authUserCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetAuthUserFromContext(ctx context.Context) *store.User {
	return ctx.Value(authUserCtx).(*store.User)
}
