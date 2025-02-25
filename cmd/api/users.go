package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/sergdort/Social/internal/store"
	"net/http"
	"strconv"
)

type userKey string

const userCtx userKey = "user"

type FollowUser struct {
	UserId int64 `json:"user_id"`
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followedUser := getUserFromContext(r)
	// TODO: Revert back to auth userID from ctx
	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.store.Follows.Follow(r.Context(), payload.UserId, followedUser.ID); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	followedUser := getUserFromContext(r)
	// TODO: Revert back to auth userID from ctx
	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.store.Follows.Unfollow(r.Context(), payload.UserId, followedUser.ID); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := getUserId(r)
		if err != nil {
			app.notFoundResponse(w, r, err)
			return
		}

		user, err := app.store.Users.GetByID(r.Context(), userId)

		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx := context.WithValue(r.Context(), userCtx, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromContext(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}

func getUserId(r *http.Request) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
}
