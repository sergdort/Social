package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrNotFound = errors.New("record not found")

const QueryTimeoutDuration = 5 * time.Second

type Storage struct {
	Posts interface {
		Create(ctx context.Context, post *Post) error
		GetByID(ctx context.Context, id int64) (*Post, error)
		Delete(ctx context.Context, id int64) error
		Update(ctx context.Context, post *Post) error
	}
	Users interface {
		Create(ctx context.Context, user *User) error
		GetByID(ctx context.Context, id int64) (*User, error)
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
