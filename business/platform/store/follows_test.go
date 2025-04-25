package store

import (
	"context"
	"errors"
	sqlc2 "github.com/sergdort/Social/business/platform/store/sqlc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestFollowsStore(t *testing.T) {

	t.Run("it should call correct query on follows", func(t *testing.T) {
		query := `-- name: CreateFollow :exec
INSERT INTO followers (user_id, follower_id)
VALUES ($1, $2)
`
		userID := int64(42)
		followerID := int64(43)
		ctx := context.Background()
		mockDB := sqlc2.NewMockDBTX(t)

		mockDB.On(
			"ExecContext",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(&FakeSqlResult{
			InsertID:      1,
			InsertError:   nil,
			AffectedRows:  0,
			AffectedError: nil,
		}, nil)

		store := FollowsStore{
			queries: sqlc2.New(mockDB),
		}

		err := store.Follow(ctx, userID, followerID)

		assert.NoError(t, err)
		mockDB.AssertCalled(t, "ExecContext", mock.Anything, query, userID, followerID)
		mockDB.AssertNumberOfCalls(t, "ExecContext", 1)
	})

	t.Run("it should return error if it fails to follow", func(t *testing.T) {
		userID := int64(42)
		followerID := int64(43)
		ctx := context.Background()
		mockDB := sqlc2.NewMockDBTX(t)

		fakeError := errors.New("something went wrong")
		mockDB.On(
			"ExecContext",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(&FakeSqlResult{
			InsertID:      0,
			InsertError:   nil,
			AffectedRows:  0,
			AffectedError: nil,
		}, fakeError)

		store := FollowsStore{
			queries: sqlc2.New(mockDB),
		}

		err := store.Follow(ctx, userID, followerID)

		assert.EqualError(t, err, fakeError.Error())
	})

	t.Run("it should call correct query on unfollow", func(t *testing.T) {
		query := `-- name: DeleteFollow :execrows
DELETE
FROM followers
WHERE user_id = $1
  AND follower_id = $2
`
		userID := int64(42)
		followerID := int64(43)
		ctx := context.Background()
		mockDB := sqlc2.NewMockDBTX(t)

		mockDB.On(
			"ExecContext",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(&FakeSqlResult{
			InsertID:      0,
			InsertError:   nil,
			AffectedRows:  1,
			AffectedError: nil,
		}, nil)

		store := FollowsStore{
			queries: sqlc2.New(mockDB),
		}

		err := store.Unfollow(ctx, userID, followerID)

		assert.NoError(t, err)
		mockDB.AssertCalled(t, "ExecContext", mock.Anything, query, userID, followerID)
		mockDB.AssertNumberOfCalls(t, "ExecContext", 1)
	})

	t.Run("it returns NotFound error if no rows affected when unfollow", func(t *testing.T) {
		userID := int64(42)
		followerID := int64(43)
		ctx := context.Background()
		mockDB := sqlc2.NewMockDBTX(t)

		mockDB.On(
			"ExecContext",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(&FakeSqlResult{
			InsertID:      0,
			InsertError:   nil,
			AffectedRows:  0,
			AffectedError: nil,
		}, nil)

		store := FollowsStore{
			queries: sqlc2.New(mockDB),
		}

		err := store.Unfollow(ctx, userID, followerID)

		assert.EqualError(t, err, ErrNotFound.Error())
	})
}
