package store

import (
	"context"
	"github.com/sergdort/Social/business/domain"
	"github.com/sergdort/Social/business/platform/store/sqlc"
	"github.com/sergdort/Social/foundation/slices"
)

type FeedStore struct {
	queries *sqlc.Queries
}

func (s *FeedStore) GetUserFeed(ctx context.Context, userId int64, q domain.PaginatedFeedQuery) ([]domain.PostWithMetadata, error) {
	feed, err := s.queries.GetUserFeed(ctx, sqlc.GetUserFeedParams{
		UserID:  userId,
		Limit:   int32(q.Limit),
		Offset:  int32(q.Offset),
		Column4: q.Search,
		Tags:    q.Tags,
	})
	if err != nil {
		return nil, err
	}
	postsWithMetadata := slices.Map(feed, convertToPostWithMetadata)
	return postsWithMetadata, nil
}
