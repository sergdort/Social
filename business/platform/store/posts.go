package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/sergdort/Social/business/domain"
	"github.com/sergdort/Social/business/platform/store/sqlc"
)

type PostStore struct {
	queries *sqlc.Queries
}

func (s *PostStore) Create(ctx context.Context, post *domain.Post) error {
	row, err := s.queries.CreatePost(ctx, sqlc.CreatePostParams{
		Content: post.Content,
		Title:   post.Title,
		UserID:  post.UserID,
		Tags:    post.Tags,
	})

	if err != nil {
		return err
	}

	post.ID = row.ID
	post.CreatedAt = row.CreatedAt.String()
	post.UpdatedAt = row.UpdatedAt.String()

	return nil
}

func (s *PostStore) GetByID(ctx context.Context, id int64) (*domain.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)

	defer cancel()

	row, err := s.queries.GetPostByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, domain.ErrNotFound
		default:
			return nil, err
		}
	}
	return &domain.Post{
		ID:        row.ID,
		Content:   row.Content,
		Title:     row.Title,
		UserID:    row.UserID,
		CreatedAt: row.CreatedAt.String(),
		UpdatedAt: row.UpdatedAt.String(),
		Tags:      row.Tags,
		Version:   int64(row.Version.Int32),
	}, nil
}

func (s *PostStore) Delete(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)

	defer cancel()

	rows, err := s.queries.DeletePostByID(ctx, id)

	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (s *PostStore) Update(ctx context.Context, post *domain.Post) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)

	defer cancel()

	version, err := s.queries.UpdatePost(ctx, sqlc.UpdatePostParams{
		Content: post.Content,
		Title:   post.Title,
		ID:      post.ID,
		Version: sql.NullInt32{
			Int32: int32(post.Version),
			Valid: true,
		},
	})

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return domain.ErrNotFound
		default:
			return err
		}
	}

	post.Version = int64(version.Int32)

	return nil
}

func convertToPostWithMetadata(feedRow sqlc.GetUserFeedRow) domain.PostWithMetadata {
	return domain.PostWithMetadata{
		Post: domain.Post{
			ID:        feedRow.ID,
			Content:   feedRow.Content,
			Title:     feedRow.Title,
			UserID:    feedRow.UserID,
			CreatedAt: feedRow.CreatedAt.String(),
			Tags:      feedRow.Tags,
			User: domain.User{
				ID:       feedRow.UserID,
				Username: feedRow.Username.String,
			},
		},
		CommentsCount: feedRow.CommentsCount,
	}
}
