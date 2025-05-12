package authapp

import (
	"context"
	"github.com/sergdort/Social/app/shared/errs"
	"github.com/sergdort/Social/business/domain"
	"github.com/sergdort/Social/foundation/jsn"
	"github.com/sergdort/Social/foundation/web"
	"net/http"
)

type authApp struct {
	useCase *domain.AuthUseCase
}

func (app *authApp) registerUserHandler(ctx context.Context, r *http.Request) web.Encoder {
	var payload domain.RegisterUserPayload
	if err := jsn.ReadJSON(r, &payload); err != nil {
		return errs.Newf(errs.InvalidArgument, "Bad Request %s", err.Error())
	}

	if err := domain.Validate.Struct(payload); err != nil {
		return errs.Newf(errs.InvalidArgument, "Bad Request %s", err.Error())
	}

	token, error := app.useCase.RegisterUser(ctx, payload)
	if error != nil {
		return errs.Newf(errs.Internal, "Failed to register user: %s", error.Error())
	}
	return token
}

// createTokenHandler godoc
//
//	@Summary		Creates a token
//	@Description	Creates a token for a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateUserTokenPayload	true	"User credentials"
//	@Success		200		{string}	string					"Token"
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/token [post]
func (app *authApp) createTokenHandler(ctx context.Context, r *http.Request) web.Encoder {
	// parse the payload creds
	var payload domain.CreateUserTokenPayload
	err := jsn.ReadJSON(r, &payload)
	// check if the user exists
	token, err := app.useCase.CreateToken(ctx, payload)
	if err != nil {
		return errs.Newf(errs.InvalidArgument, "Invalid email or password")
	}
	return TokenResponse{Token: token}
}
