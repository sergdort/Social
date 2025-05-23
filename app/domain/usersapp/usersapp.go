package usersapp

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/sergdort/Social/app/shared/errs"
	"github.com/sergdort/Social/app/shared/mid"
	"github.com/sergdort/Social/business/domain"
	"github.com/sergdort/Social/foundation/web"
)

type userKey string

const userCtx userKey = "user"

type authUserKey string

const authUserCtx userKey = "authUser"

type FollowUser struct {
	UserId int64 `json:"user_id"`
}

type userApp struct {
	usersUseCase *domain.UsersUseCase
}

func newApp(usersUseCase *domain.UsersUseCase) *userApp {
	return &userApp{
		usersUseCase: usersUseCase,
	}
}

// GetUser godoc
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	User
//	@Failure		400	{object}	error
//	@Failure		204	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{id} [get]
func (app *userApp) getUserHandler(ctx context.Context, r *http.Request) web.Encoder {
	user, err := getUserFromContext(ctx)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	return web.Response[User]{Data: toAppUser(user)}
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
func (app *userApp) followUserHandler(ctx context.Context, r *http.Request) web.Encoder {
	currentUserID, err := mid.GetAuthUserID(ctx)

	if err != nil {
		return errs.New(errs.Internal, err)
	}

	userToFollow, err := getUserFromContext(ctx)

	if err != nil {
		return errs.New(errs.Internal, err)
	}

	err = app.usersUseCase.FollowUser(ctx, currentUserID, userToFollow.ID)

	if err != nil {
		return errs.New(errs.Internal, err)
	}
	return web.NewNoResponse()
}

func (app *userApp) unfollowUserHandler(ctx context.Context, r *http.Request) web.Encoder {
	currentUserID, err := mid.GetAuthUserID(ctx)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	userToUnfollow, err := getUserFromContext(ctx)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	err = app.usersUseCase.UnfollowUser(ctx, currentUserID, userToUnfollow.ID)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	return web.NewNoResponse()
}

func (app *userApp) activateUserHandler(ctx context.Context, r *http.Request) web.Encoder {
	token := web.Param(r, "token")

	err := app.usersUseCase.ActivateUser(ctx, token)
	if err != nil {
		return errs.New(errs.Internal, err)
	}
	return web.NewNoResponse()
}

func (app *userApp) userContextMiddleware(useCase *domain.UsersUseCase) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			userId, err := getUserId(r)
			if err != nil {
				return errs.Newf(errs.Internal, "qury by id: %s", err)
			}
			user, err := useCase.GetUserById(ctx, userId)
			if err != nil {
				return errs.Newf(errs.NotFound, "User not found")
			}
			return next(context.WithValue(ctx, userCtx, user), r)
		}
		return h
	}
	return m
}

func getAuthUserFromContext(r *http.Request) *domain.User {
	return r.Context().Value(authUserCtx).(*domain.User)
}

func getUserFromContext(context context.Context) (*domain.User, error) {
	user, ok := context.Value(userCtx).(*domain.User)
	if !ok {
		return nil, errors.New("user not found in context")
	}
	return user, nil
}

func getUserId(r *http.Request) (int64, error) {
	return strconv.ParseInt(web.Param(r, "userID"), 10, 64)
}
