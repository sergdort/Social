package main

import (
	"github.com/sergdort/Social/business/domain"
	"net/http"
)

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
//	@Success		200		{object}	[]domain.PostWithMetadata
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	query := domain.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		SortBy: "desc",
		Tags:   []string{},
	}
	ParsePaginatedFeedQuery(&query, r)

	if err := Validate.Struct(query); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	ctx := r.Context()
	user := getAuthUserFromContext(r)
	posts, err := app.store.Posts.GetUserFeed(ctx, user.ID, query)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, posts); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
