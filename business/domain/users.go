package domain

import (
	"context"
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  Password `json:"-"`
	CreatedAt string   `json:"created_at"`
	IsActive  bool     `json:"is_active"`
	RoleID    int64    `json:"role_id"`
	Role      Role     `json:"role"`
}

type Password struct {
	Text *string
	Hash []byte
}

func (p *Password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.Hash = hash
	p.Text = &text

	return nil
}

func (p *Password) Verify(text string) error {
	return bcrypt.CompareHashAndPassword(p.Hash, []byte(text))
}

type UsersUseCase struct {
	cache       UsersCache
	usersRepo   UsersRepository
	followsRepo FollowsRepository
}

func NewUsersUseCase(
	cache UsersCache,
	usersRepo UsersRepository,
	followsRepo FollowsRepository,
) *UsersUseCase {
	return &UsersUseCase{
		cache:       cache,
		usersRepo:   usersRepo,
		followsRepo: followsRepo,
	}
}

func (uc *UsersUseCase) GetUserById(ctx context.Context, userID int64) (*User, error) {
	// Try to get user from cache
	if user, err := uc.cache.Get(ctx, userID); err == nil && user != nil {
		return user, nil
	}

	// Fetch from database
	user, err := uc.usersRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Try to store in cache but don't fail if caching fails
	if err := uc.cache.Set(ctx, user); err != nil {
		// TODO: Log error?
	}
	return user, nil
}

func (uc *UsersUseCase) FollowUser(ctx context.Context, userID int64, followerID int64) error {
	return uc.followsRepo.Follow(ctx, userID, followerID)
}

func (uc *UsersUseCase) UnfollowUser(ctx context.Context, userID int64, followerID int64) error {
	return uc.followsRepo.Unfollow(ctx, userID, followerID)
}

func (uc *UsersUseCase) ActivateUser(ctx context.Context, token string) error {
	return uc.usersRepo.Activate(ctx, token)
}

type UsersCache interface {
	Get(ctx context.Context, id int64) (*User, error)
	Set(ctx context.Context, user *User) error
}

type UsersRepository interface {
	Create(ctx context.Context, tx *sql.Tx, user *User) error
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	CreateAndInvite(ctx context.Context, user *User, token string, expiration time.Duration) error
	RevertCreateAndInvite(ctx context.Context, id int64) error
	Activate(ctx context.Context, token string) error
}
