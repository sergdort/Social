package mid

import (
	"context"
	"errors"
	"github.com/sergdort/Social/business/domain"
)

type ctxKey int

const (
	userIDKey = iota + 1
	userKey
)

func setAuthUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetAuthUserID returns the user id from the context.
func GetAuthUserID(ctx context.Context) (int64, error) {
	v, ok := ctx.Value(userIDKey).(int64)
	if !ok {
		return 0, errors.New("user id not found in context")
	}

	return v, nil
}

func setUser(ctx context.Context, usr domain.User) context.Context {
	return context.WithValue(ctx, userKey, usr)
}

// GetUser returns the user from the context.
func GetUser(ctx context.Context) (domain.User, error) {
	v, ok := ctx.Value(userKey).(domain.User)
	if !ok {
		return domain.User{}, errors.New("user not found in context")
	}

	return v, nil
}
