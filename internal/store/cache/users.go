package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sergdort/Social/internal/store"
	"time"
)

type UsersStore struct {
	rdb *redis.Client
}

func (s *UsersStore) Get(ctx context.Context, id int64) (*store.User, error) {
	key := cacheKey(id)
	data, err := s.rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var user store.User
	if data == "" {
		return nil, nil
	}
	err = json.Unmarshal([]byte(data), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UsersStore) Set(ctx context.Context, user *store.User) error {
	key := cacheKey(user.ID)

	json, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = s.rdb.SetEx(ctx, key, string(json), time.Hour).Err()
	return err
}

func cacheKey(userID int64) string {
	return fmt.Sprintf("user-%d", userID)
}
