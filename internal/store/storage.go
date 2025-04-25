package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/sergdort/Social/business/domain"
	"github.com/sergdort/Social/internal/store/sqlc"
	"time"
)

var ErrNotFound = errors.New("record not found")
var ErrDuplicateEmail = errors.New("email already exists")
var ErrDuplicateUsername = errors.New("username already exists")

const QueryTimeoutDuration = 5 * time.Second

type Storage struct {
	Posts    PostsRepository
	Users    domain.UsersRepository
	Comments CommentsRepository
	Follows  FollowsRepository
	Roles    domain.RolesRepository
}

type PostsRepository interface {
	Create(ctx context.Context, post *Post) error
	GetByID(ctx context.Context, id int64) (*Post, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, post *Post) error
	GetUserFeed(ctx context.Context, userId int64, query PaginatedFeedQuery) ([]PostWithMetadata, error)
}

type CommentsRepository interface {
	Create(ctx context.Context, comment *Comment) error
	GetAllByPostID(ctx context.Context, postID int64) ([]Comment, error)
}

type FollowsRepository interface {
	Follow(ctx context.Context, userID int64, followerID int64) error
	Unfollow(ctx context.Context, userID int64, followerID int64) error
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{sqlc.New(db)},
		Users:    &UserStore{db, sqlc.New(db)},
		Comments: &CommentStore{sqlc.New(db)},
		Follows:  &FollowsStore{sqlc.New(db)},
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
