package main

import (
	"github.com/sergdort/Social/internal/store"
	"net/http"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: use auth user to fetch feed for
	query := store.PaginatedFeedQuery{
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

	posts, err := app.store.Posts.GetUserFeed(r.Context(), int64(314), query)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, posts); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
