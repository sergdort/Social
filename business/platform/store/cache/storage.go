package cache

import (
	"github.com/redis/go-redis/v9"
	"github.com/sergdort/Social/business/domain"
)

type Storage struct {
	Users domain.UsersCache
}

func NewStorage(rdb *redis.Client) Storage {
	return Storage{
		Users: &UsersStore{rdb: rdb},
	}
}
