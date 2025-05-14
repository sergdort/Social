package domain

import "context"

type PaginatedFeedQuery struct {
	Limit  int      `json:"limit" validate:"gte=1,lte=50"`
	Offset int      `json:"offset" validate:"gte=0"`
	SortBy string   `form:"sort_by" validate:"oneof=asc desc"`
	Tags   []string `json:"tags" validate:"max=5"`
	Search string   `form:"search" validate:"max=100"`
	Since  string   `json:"since"`
	Until  string   `json:"until"`
}

type FeedUseCase interface {
	GetUserFeed(ctx context.Context, userId int64, query PaginatedFeedQuery) ([]PostWithMetadata, error)
}
