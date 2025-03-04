package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrNotFound = errors.New("record not found")
var ErrDuplicateEmail = errors.New("email already exists")
var ErrDuplicateUsername = errors.New("username already exists")

const QueryTimeoutDuration = 5 * time.Second

type Storage struct {
	Posts interface {
		Create(ctx context.Context, post *Post) error
		GetByID(ctx context.Context, id int64) (*Post, error)
		Delete(ctx context.Context, id int64) error
		Update(ctx context.Context, post *Post) error
		GetUserFeed(ctx context.Context, userId int64, query PaginatedFeedQuery) ([]PostWithMetadata, error)
	}
	Users interface {
		Create(ctx context.Context, tx *sql.Tx, user *User) error
		GetByID(ctx context.Context, id int64) (*User, error)
		CreateAndInvite(ctx context.Context, user *User, token string, expiration time.Duration) error
	}
	Comments interface {
		Create(ctx context.Context, comment *Comment) error
		GetAllByPostID(ctx context.Context, postID int64) ([]Comment, error)
	}
	Follows interface {
		Follow(ctx context.Context, userID int64, followerID int64) error
		Unfollow(ctx context.Context, userID int64, followerID int64) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{db},
		Users:    &UserStore{db},
		Comments: &CommentStore{db},
		Follows:  &FollowsStore{db},
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
