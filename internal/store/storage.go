package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/sergdort/Social/internal/store/sqlc"
	"time"
)

var ErrNotFound = errors.New("record not found")
var ErrDuplicateEmail = errors.New("email already exists")
var ErrDuplicateUsername = errors.New("username already exists")

const QueryTimeoutDuration = 5 * time.Second

type Storage struct {
	Posts    PostsRepository
	Users    UsersRepository
	Comments CommentsRepository
	Follows  FollowsRepository
	Roles    RolesRepository
}

type PostsRepository interface {
	Create(ctx context.Context, post *Post) error
	GetByID(ctx context.Context, id int64) (*Post, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, post *Post) error
	GetUserFeed(ctx context.Context, userId int64, query PaginatedFeedQuery) ([]PostWithMetadata, error)
}

type UsersRepository interface {
	Create(ctx context.Context, tx *sql.Tx, user *User) error
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	CreateAndInvite(ctx context.Context, user *User, token string, expiration time.Duration) error
	RevertCreateAndInvite(ctx context.Context, id int64) error
	Activate(ctx context.Context, token string) error
}

type CommentsRepository interface {
	Create(ctx context.Context, comment *Comment) error
	GetAllByPostID(ctx context.Context, postID int64) ([]Comment, error)
}

type FollowsRepository interface {
	Follow(ctx context.Context, userID int64, followerID int64) error
	Unfollow(ctx context.Context, userID int64, followerID int64) error
}

type RolesRepository interface {
	GetByRoleType(ctx context.Context, name RoleType) (*Role, error)
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{db, sqlc.New(db)},
		Users:    &UserStore{db, sqlc.New(db)},
		Comments: &CommentStore{db},
		Follows:  &FollowsStore{db},
		Roles:    &RolesStore{queries: sqlc.New(db)},
	}
}

func withTransaction(db *sql.DB, ctx context.Context, f func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := f(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
