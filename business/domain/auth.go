package domain

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"time"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

type RegisterUserPayload struct {
	UserName string `json:"user_name" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

type InvitationToken struct {
	Token         string `json:"token"`
	InvitationURL string `json:"invitation_url"`
}

func (token InvitationToken) Encode() (data []byte, contentType string, err error) {
	jsonData, err := json.Marshal(token)
	return jsonData, "application/json", err
}

type CreateUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

type AuthConfig struct {
	InvitationExp time.Duration
	FrontendURL   string
}

type Claims struct {
	UserID int64
}

type TokenGenerator interface {
	GenerateToken(ctx context.Context, userID int64) (string, error)
}

type TokenValidator interface {
	ValidateToken(ctx context.Context, token string) (Claims, error)
}

type AuthUseCase struct {
	config     AuthConfig
	roles      RolesRepository
	users      UsersRepository
	token      TokenGenerator
	tokenValid TokenValidator
}

func NewAuthUseCase(config AuthConfig, roles RolesRepository, users UsersRepository, token TokenGenerator, tokenValid TokenValidator) *AuthUseCase {
	return &AuthUseCase{
		config:     config,
		roles:      roles,
		users:      users,
		token:      token,
		tokenValid: tokenValid,
	}
}

func (auth *AuthUseCase) RegisterUser(ctx context.Context, payload RegisterUserPayload) (*InvitationToken, error) {
	role, err := auth.roles.GetByRoleType(ctx, RoleTypeUser)
	if err != nil {
		return nil, err
	}

	user := &User{
		Username: payload.UserName,
		Email:    payload.Email,
		RoleID:   role.ID,
	}

	if err := user.Password.Set(payload.Password); err != nil {
		return nil, err
	}

	token := uuid.New().String()
	response := InvitationToken{
		Token:         token,
		InvitationURL: fmt.Sprintf("%s/confirm/%s", auth.config.FrontendURL, token),
	}

	hashToken := hashToken(response.Token)

	if err := auth.users.CreateAndInvite(ctx, user, hashToken, auth.config.InvitationExp); err != nil {
		return nil, err
	}

	return &response, nil
}

func (auth *AuthUseCase) ActivateToken(ctx context.Context, token string) error {
	return auth.users.Activate(ctx, hashToken(token))
}

func (auth *AuthUseCase) CreateToken(ctx context.Context, payload CreateUserTokenPayload) (string, error) {
	user, err := auth.users.GetByEmail(ctx, payload.Email)
	if err != nil {
		return "", err
	}
	if err := user.Password.Verify(payload.Password); err != nil {
		return "", err
	}
	return auth.token.GenerateToken(ctx, user.ID)
}

func (auth *AuthUseCase) ValidateToken(ctx context.Context, token string) (Claims, error) {
	return auth.tokenValid.ValidateToken(ctx, token)
}

func hashToken(plainToken string) string {
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])
	return hashToken
}
