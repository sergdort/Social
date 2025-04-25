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
	Posts    domain.PostsRepository
	Users    domain.UsersRepository
	Comments domain.CommentsRepository
	Follows  domain.FollowsRepository
	Roles    domain.RolesRepository
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
