package store

import (
	"context"
	"database/sql"
	"github.com/sergdort/Social/business/domain"
	"github.com/sergdort/Social/business/platform/store/sqlc"
	"time"
)

const QueryTimeoutDuration = 5 * time.Second

type Storage struct {
	Posts    domain.PostsRepository
	Users    domain.UsersRepository
	Comments domain.CommentsRepository
	Follows  domain.FollowsRepository
	Roles    domain.RolesRepository
	Feed     domain.FeedRepository
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{sqlc.New(db)},
		Users:    &UserStore{db, sqlc.New(db)},
		Comments: &CommentStore{sqlc.New(db)},
		Follows:  &FollowsStore{sqlc.New(db)},
		Roles:    &RolesStore{queries: sqlc.New(db)},
		Feed:     &FeedStore{sqlc.New(db)},
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
