package store

import (
	"context"
	"database/sql"
	"github.com/sergdort/Social/business/domain"
	sqlc2 "github.com/sergdort/Social/business/platform/store/sqlc"
	"github.com/sergdort/Social/foundation/slices"
)

type CommentStore struct {
	queries *sqlc2.Queries
}

func (s *CommentStore) Create(ctx context.Context, comment *domain.Comment) error {

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)

	defer cancel()

	result, err := s.queries.CreateComment(ctx, sqlc2.CreateCommentParams{
		PostID: comment.PostID,
		UserID: comment.UserID,
		Content: sql.NullString{
			String: comment.Content,
			Valid:  true,
		},
	})

	if err != nil {
		return err
	}

	comment.ID = result.ID
	comment.CreatedAt = result.CreatedAt.String()

	return nil
}

func (s *CommentStore) GetAllByPostID(ctx context.Context, postID int64) ([]domain.Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)

	defer cancel()

	rows, err := s.queries.GetAllCommentsByPostID(ctx, postID)
	if err != nil {
		return nil, err
	}

	comments := slices.Map(rows, convertToComment)

	return comments, nil
}

func convertToComment(row sqlc2.GetAllCommentsByPostIDRow) domain.Comment {
	return domain.Comment{
		ID:        row.ID,
		PostID:    row.PostID,
		UserID:    row.UserID,
		Content:   row.Content.String,
		CreatedAt: row.CreatedAt.String(),
		User: domain.User{
			ID:       row.UserID,
			Username: row.Username,
		},
	}
}
