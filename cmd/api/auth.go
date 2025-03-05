package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sergdort/Social/internal/mailer"
	"github.com/sergdort/Social/internal/store"
	"net/http"
)

type RegisterUserPayload struct {
	UserName string `json:"user_name" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

type InvitationTokenResponse struct {
	Token string `json:"token"`
}

// registerUserHandler godoc
//
//	@Summary		Registers a user
//	@Description	Registers a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User credentials"
//	@Success		201		{object}	InvitationTokenResponse			"User registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &store.User{
		Username: payload.UserName,
		Email:    payload.Email,
	}

	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	response := InvitationTokenResponse{
		Token: uuid.New().String(),
	}
	hashToken := hashToken(response.Token)

	if err := app.store.Users.CreateAndInvite(ctx, user, hashToken, app.config.mail.exp); err != nil {
		switch {
		case errors.Is(err, store.ErrDuplicateEmail):
			app.badRequestResponse(w, r, err)
		case errors.Is(err, store.ErrDuplicateUsername):
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	activationURL := fmt.Sprintf("%s/confirm/%s", app.config.frontEndURL, response.Token)
	fmt.Println(
		"activationURL",
		activationURL,
	)

	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}
	if err := app.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars); err != nil {
		app.internalServerError(w, r, err)
		app.logger.Errorw("Error sending welcome email: ", "error", err.Error(), "vars", vars)

		// rollback user creation
		if err := app.store.Users.RevertCreateAndInvite(ctx, user.ID); err != nil {
			app.logger.Errorw("Error rolling back user creation: ", "error", err.Error(), "userID", user.ID)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, response); err != nil {
		app.internalServerError(w, r, err)
	}
}

// ActivateUser godoc
//
//	@Summary		Activates/Register a user
//	@Description	Activates/Register a user by invitation token
//	@Tags			users
//	@Produce		json
//	@Param			token	path		string	true	"Invitation token"
//	@Success		204		{string}	string	"User activated"
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/activate/{token} [put]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	err := app.store.Users.Activate(r.Context(), hashToken(token))

	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func hashToken(plainToken string) string {
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])
	return hashToken
}
