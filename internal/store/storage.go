package store

import (
	"context"
	"database/sql"
	"errors"
)

var ErrNotFound = errors.New("record not found")

type Storage struct {
	Posts interface {
		Create(ctx context.Context, post *Post) error
		GetPostByID(ctx context.Context, id int64) (*Post, error)
		Delete(ctx context.Context, id int64) error
		Update(ctx context.Context, post *Post) error
	}
	Users interface {
		Create(ctx context.Context, user *User) error
	}
	Comments interface {
		GetAllByPostID(ctx context.Context, postID int64) ([]Comment, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{db},
		Users:    &UserStore{db},
		Comments: &CommentStore{db},
	}
}
