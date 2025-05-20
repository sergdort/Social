package store

import (
	"context"
	"github.com/sergdort/Social/business/domain"
	sqlc2 "github.com/sergdort/Social/business/platform/store/sqlc"
)

type FollowsStore struct {
	queries *sqlc2.Queries
}

func (s *FollowsStore) Follow(ctx context.Context, userId int64, followerId int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeoutDuration)
	defer cancel()

	err := s.queries.CreateFollow(ctx, sqlc2.CreateFollowParams{
		UserID:     userId,
		FollowerID: followerId,
	})

	return err
}

func (s *FollowsStore) Unfollow(ctx context.Context, userID int64, followerID int64) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration) // Use passed ctx
	defer cancel()

	rows, err := s.queries.DeleteFollow(ctx, sqlc2.DeleteFollowParams{
		UserID:     userID,
		FollowerID: followerID,
	})
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}
