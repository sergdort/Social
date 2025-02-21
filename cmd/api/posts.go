package main

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/sergdort/Social/internal/store"
	"golang.org/x/net/context"
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

func (app *application) createPostsHandler(w http.ResponseWriter, r *http.Request) {
	var userId = 1
	var ctx = r.Context()
	var payload CreatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
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

	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromContext(r)

	var comments, err3 = app.store.Comments.GetAllByPostID(r.Context(), post.ID)
	if err3 != nil {
		app.internalServerError(w, r, err3)
		return
	}

	post.Comments = comments

	_ = app.jsonResponse(w, http.StatusOK, post)
}

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
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

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
		case errors.Is(err, store.ErrNotFound):
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

		var post, err2 = app.store.Posts.GetPostByID(r.Context(), postID)

		if err2 != nil {
			switch {
			case errors.Is(err2, store.ErrNotFound):
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

func getPostFromContext(r *http.Request) *store.Post {
	post, _ := r.Context().Value(postCtx).(*store.Post)
	return post
}

func (p *UpdatePostPayload) update(post *store.Post) {
	if p.Title != nil {
		post.Title = *p.Title
	}
	if p.Content != nil {
		post.Content = *p.Content
	}
}
