package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/sergdort/Social/business/domain"
	"github.com/sergdort/Social/business/platform/store"
	"net/http"
	"strconv"
)

type userKey string

const userCtx userKey = "user"

type FollowUser struct {
	UserId int64 `json:"user_id"`
}

// FollowUser godoc
//
//	@Summary		Follows a user
//	@Description	Follows a user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		204	{string}	No	Content
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{id}/follow [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	currentUser := getAuthUserFromContext(r)
	userToFollow := getUserFromContext(r)

	if err := app.store.Follows.Follow(ctx, currentUser.ID, userToFollow.ID); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// UnfollowUser godoc
//
//	@Summary		Unfollows a user
//	@Description	Unfollows a user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		204	{string}	No	Content
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{id}/unfollow [put]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	currentUser := getAuthUserFromContext(r)
	userToUnfollow := getUserFromContext(r)

	if err := app.store.Follows.Unfollow(ctx, currentUser.ID, userToUnfollow.ID); err != nil {
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

func getUserFromContext(r *http.Request) *domain.User {
	user, _ := r.Context().Value(userCtx).(*domain.User)
	return user
}

func getUserId(r *http.Request) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
}
