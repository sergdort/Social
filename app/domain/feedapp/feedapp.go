package feedapp

import (
	"context"
	"net/http"

	"github.com/sergdort/Social/app/shared/mid"

	"github.com/sergdort/Social/app/shared/errs"
	"github.com/sergdort/Social/business/domain"
	"github.com/sergdort/Social/foundation/slices"
	"github.com/sergdort/Social/foundation/web"
)

type feedApp struct {
	feedUseCase domain.FeedRepository
}

// getUserFeedHandler godoc
//
//	@Summary		Fetches the user feed
//	@Description	Fetches the user feed
//	@Tags			feed
//	@Accept			json
//	@Produce		json
//	@Param			since	query		string	false	"Since"
//	@Param			until	query		string	false	"Until"
//	@Param			limit	query		int		false	"Limit"
//	@Param			offset	query		int		false	"Offset"
//	@Param			sort_by	query		string	false	"Sort"
//	@Param			tags	query		string	false	"Tags"
//	@Param			search	query		string	false	"Search"
//	@Success		200		{object}	FeedData
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/feed [get]
func (app *feedApp) getFeedHandler(ctx context.Context, r *http.Request) web.Encoder {
	query := domain.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		SortBy: "desc",
		Tags:   []string{},
	}
	parsePaginatedFeedQuery(&query, r)

	if err := domain.Validate.Struct(query); err != nil {
		return errs.Newf(errs.InvalidArgument, err.Error())
	}
	userID, err := mid.GetAuthUserID(ctx)
	if err != nil {
		return errs.Newf(errs.Internal, "feed query could not get user id %s", err.Error())
	}

	feed, err := app.feedUseCase.GetUserFeed(ctx, userID, query)
	if err != nil {
		return errs.Newf(errs.Internal, "could not get user feed %s", err.Error())
	}

	feedItems := slices.Map(feed, toPostFeedItem)
	return web.NewResponse(feedItems)
}
