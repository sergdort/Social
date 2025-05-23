package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/sergdort/Social/business/domain"
	"net/http"
	"strconv"
)

type postKey string

const postCtx postKey = "post"

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}

// GetPost godoc
//
//	@Summary		Fetches a post
//	@Description	Fetches a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Post ID"
//	@Success		200	{object}	domain.Post
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [get]
func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromContext(r)

	var comments, err3 = app.store.Comments.GetAllByPostID(r.Context(), post.ID)
	if err3 != nil {
		app.internalServerError(w, r, err3)
		return
	}

	post.Comments = comments

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

// DeletePost godoc
//
//	@Summary		Deletes a post
//	@Description	RevertCreateAndInvite a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Post ID"
//	@Success		204	{object} string
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [delete]
func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	var postID, err = app.getPostId(r)
	if err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	var ctx = r.Context()

	err = app.store.Posts.Delete(ctx, postID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdatePost godoc
//
//	@Summary		Updates a post
//	@Description	Updates a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Post ID"
//	@Param			payload	body		UpdatePostPayload	true	"Post payload"
//	@Success		200		{object}	domain.Post
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [patch]
func (app *application) patchPostsHandler(w http.ResponseWriter, r *http.Request) {
	var payload = UpdatePostPayload{}
	var post = getPostFromContext(r)

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	payload.update(post)

	err := app.store.Posts.Update(r.Context(), post)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postID, err := app.getPostId(r)
		if err != nil {
			_ = writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		var post, err2 = app.store.Posts.GetByID(r.Context(), postID)

		if err2 != nil {
			switch {
			case errors.Is(err2, domain.ErrNotFound):
				app.notFoundResponse(w, r, err2)
			default:
				app.internalServerError(w, r, err2)
			}
			return
		}

		ctx := context.WithValue(r.Context(), postCtx, post)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) getPostId(r *http.Request) (int64, error) {
	var idParam = chi.URLParam(r, "postID")
	var postID, err = strconv.ParseInt(idParam, 10, 64)
	return postID, err
}

func getPostFromContext(r *http.Request) *domain.Post {
	post, _ := r.Context().Value(postCtx).(*domain.Post)
	return post
}

func (p *UpdatePostPayload) update(post *domain.Post) {
	if p.Title != nil {
		post.Title = *p.Title
	}
	if p.Content != nil {
		post.Content = *p.Content
	}
}
