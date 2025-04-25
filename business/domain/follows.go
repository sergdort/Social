package domain

import "context"

type FollowsRepository interface {
	Follow(ctx context.Context, userID int64, followerID int64) error
	Unfollow(ctx context.Context, userID int64, followerID int64) error
}
