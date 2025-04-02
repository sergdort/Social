package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sergdort/Social/internal/mailer"
	"github.com/sergdort/Social/internal/store"
	"math/rand"
	"net/http"
	"time"
)

type RegisterUserPayload struct {
	UserName string `json:"user_name" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

type InvitationTokenResponse struct {
	Token string `json:"token"`
}

type CreateUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

// registerUserHandler godoc
//
//	@Summary		Registers a user
//	@Description	Registers a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload		true	"User credentials"
//	@Success		201		{object}	InvitationTokenResponse	"User registered"
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

	role, err := app.store.Roles.GetByRoleType(r.Context(), store.RoleTypeUser)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	user := &store.User{
		Username: payload.UserName,
		Email:    payload.Email,
		RoleID:   role.ID,
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

	go func() {
		_ = retryRevert(4, func() error {
			return app.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars)
		}, func() error {
			return app.store.Users.RevertCreateAndInvite(ctx, user.ID)
		})
	}()

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
func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	// parse the payload creds
	var payload CreateUserTokenPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// check if the user exists
	user, err := app.store.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("invalid email or password"))
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	if err := user.Password.Verify(payload.Password); err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("invalid email or password"))
		return
	}
	// generate the token -> add claims

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(app.config.auth.jwt.exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.config.auth.jwt.tokenHost,
		"aud": app.config.auth.jwt.tokenHost,
	}
	token, err := app.authenticator.GenerateToken(claims)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// send it to the client

	if err := app.jsonResponse(w, http.StatusCreated, struct {
		Token string `json:"token"`
	}{
		Token: token,
	}); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func hashToken(plainToken string) string {
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])
	return hashToken
}

func retryRevert(times int, retry func() error, revert func() error) error {
	for i := 0; i < times; i++ {
		err := retry()
		if err != nil {
			exponentialBackoffWithJitter(i)
			continue
		}
		return nil
	}
	return revert()
}

func exponentialBackoffWithJitter(retries int) {
	for i := 0; i < retries; i++ {
		sleepTime := time.Second * time.Duration(1<<i)              // 2^i seconds
		jitter := time.Duration(rand.Intn(1000)) * time.Millisecond // Add up to 1s jitter
		time.Sleep(sleepTime + jitter)
	}
}
