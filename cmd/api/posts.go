package main

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/sergdort/Social/internal/store"
	"net/http"
	"strconv"
)

type CreatePostPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostsHandler(w http.ResponseWriter, r *http.Request) {
	var userId = 1
	var ctx = r.Context()
	var payload CreatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var post = &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		// TODO: change after auth
		UserID: int64(userId),
	}

	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	var idParam = chi.URLParam(r, "postID")
	var postID, err = strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	var ctx = r.Context()

	var post, err2 = app.store.Posts.GetPostByID(ctx, int64(postID))
	if err2 != nil {
		switch {
		case errors.Is(err2, store.ErrNotFound):
			app.notFoundResponse(w, r, err2)
		default:
			app.internalServerError(w, r, err2)
		}
		return
	}
	_ = writeJSON(w, http.StatusOK, post)
}
