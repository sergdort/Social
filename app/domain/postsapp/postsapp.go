package postsapp

import (
	"context"
	"errors"
	"github.com/sergdort/Social/app/shared/errs"
	"github.com/sergdort/Social/app/shared/mid"
	"github.com/sergdort/Social/business/domain"
	"github.com/sergdort/Social/foundation/jsn"
	"github.com/sergdort/Social/foundation/web"
	"net/http"
	"strconv"
)

type postsApp struct {
	repo domain.PostsRepository
}

type postKey string

const postCtx postKey = "post"

// CreatePost godoc
//
//	@Summary		Creates a post
//	@Description	Creates a post
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreatePostPayload	true	"Post Payload"
//	@Success		201		{object}	domain.Post
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/ [post]
func (app *postsApp) createPostsHandler(ctx context.Context, r *http.Request) web.Encoder {
	var payload CreatePostPayload

	if err := jsn.ReadJSON(r, &payload); err != nil {
		return errs.Newf(errs.InvalidArgument, "invalid payload %s", err.Error())
	}

	if err := domain.Validate.Struct(payload); err != nil {
		return errs.Newf(errs.InvalidArgument, err.Error())
	}

	userID, err := mid.GetAuthUserID(ctx)

	if err != nil {
		return errs.Newf(errs.Internal, err.Error())
	}

	var post = &domain.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  userID,
	}

	if err := app.repo.Create(ctx, post); err != nil {
		return errs.Newf(errs.Internal, err.Error())
	}
	return web.NewResponse(post)
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
func (app *postsApp) getPostHandler(ctx context.Context, r *http.Request) web.Encoder {
	post, err := getPostFromContext(ctx)

	if err != nil {
		return errs.Newf(errs.Internal, err.Error())
	}

	return web.NewResponse(post)
}

func (app *postsApp) postsContextMiddleware() web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			postId, err := strconv.ParseInt(web.Param(r, "postId"), 10, 64)
			if err != nil {
				return errs.Newf(errs.InvalidArgument, "invalid postId %s", err.Error())
			}
			post, err := app.repo.GetByID(ctx, postId)
			if err != nil {
				switch {
				case errors.Is(err, domain.ErrNotFound):
					return errs.New(errs.NotFound, domain.ErrNotFound)
				default:
					return errs.New(errs.Internal, err)
				}
			}
			return next(context.WithValue(ctx, postCtx, post), r)
		}
		return h
	}
	return m
}

func getPostFromContext(ctx context.Context) (*domain.Post, error) {
	post, ok := ctx.Value(postCtx).(*domain.Post)
	if !ok {
		return nil, errors.New("No post found in context")
	}
	return post, nil
}
