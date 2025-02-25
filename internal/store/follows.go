package store

import (
	"context"
	"database/sql"
)

type FollowsStore struct {
	db *sql.DB
}

func (s *FollowsStore) Follow(ctx context.Context, userId int64, followerId int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(
		ctx,
		`INSERT INTO followers (user_id, follower_id) VALUES ($1, $2)`,
		userId,
		followerId,
	)

	return err
}

func (s *FollowsStore) Unfollow(ctx context.Context, userID int64, followerID int64) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration) // Use passed ctx
	defer cancel()

	res, err := s.db.ExecContext(
		ctx,
		`DELETE FROM followers WHERE user_id = $1 AND follower_id = $2`,
		userID,
		followerID,
	)
	if err != nil {
		return err
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrNotFound
	}

	return nil
}
