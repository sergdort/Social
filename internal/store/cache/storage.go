package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/sergdort/Social/internal/store"
)

type Storage struct {
	Users interface {
		Get(ctx context.Context, id int64) (*store.User, error)
		Set(ctx context.Context, user *store.User) error
	}
}

func NewStorage(rdb *redis.Client) Storage {
	return Storage{
		Users: &UsersStore{rdb: rdb},
	}
}
