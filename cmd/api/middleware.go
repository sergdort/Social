package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/sergdort/Social/business/domain"
	"net/http"
	"strings"
)

type authUserKey string

const authUserCtx authUserKey = "authUser"

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

func (app *application) checkPostOwnershipMiddleware(roleType domain.RoleType, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currentUser := getAuthUserFromContext(r)
		post := getPostFromContext(r)

		if currentUser.ID == post.UserID {
			next.ServeHTTP(w, r)
			return
		}

		// Role precedence check
		isValid, err := app.checkRolePrecedenceForPost(currentUser, r.Context(), roleType)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}
		if !isValid {
			app.forbiddenResponse(w, r, fmt.Errorf("forbidden"))
		}

		next.ServeHTTP(w, r)
	})
}

func getAuthUserFromContext(r *http.Request) *domain.User {
	return r.Context().Value(authUserCtx).(*domain.User)
}

func (app *application) checkRolePrecedenceForPost(user *domain.User, ctx context.Context, roleType domain.RoleType) (bool, error) {
	role, err := app.store.Roles.GetByRoleType(ctx, roleType)
	if err != nil {
		return false, err
	}
	return user.Role.Level >= role.Level, nil
}

func (app *application) getUser(ctx context.Context, userID int64) (*domain.User, error) {
	// Try to get user from cache
	if user, err := app.cache.Users.Get(ctx, userID); err == nil && user != nil {
		return user, nil
	}

	// Fetch from database
	user, err := app.store.Users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Try to store in cache but don't fail if caching fails
	if err := app.cache.Users.Set(ctx, user); err != nil {
		app.logger.Warn(ctx, "Failed to set user in cache", "userID", userID, "error", err)
	}

	return user, nil
}
