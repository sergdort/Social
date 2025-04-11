package store

import (
	"context"
	"database/sql"
	"github.com/sergdort/Social/internal/store/sqlc"
)

type Comment struct {
	ID        int64  `json:"id"`
	PostID    int64  `json:"post_id"`
	UserID    int64  `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`

	User User `json:"user"`
}

type CommentStore struct {
	queries *sqlc.Queries
}

func (s *CommentStore) Create(ctx context.Context, comment *Comment) error {

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)

	defer cancel()

	result, err := s.queries.CreateComment(ctx, sqlc.CreateCommentParams{
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

func (s *CommentStore) GetAllByPostID(ctx context.Context, postID int64) ([]Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)

	defer cancel()

	rows, err := s.queries.GetAllCommentsByPostID(ctx, postID)
	if err != nil {
		return nil, err
	}

	comments := Map(rows, convertToComment)

	return comments, nil
}

func convertToComment(row sqlc.GetAllCommentsByPostIDRow) Comment {
	return Comment{
		ID:        row.ID,
		PostID:    row.PostID,
		UserID:    row.UserID,
		Content:   row.Content.String,
		CreatedAt: row.CreatedAt.String(),
		User: User{
			ID:       row.UserID,
			Username: row.Username,
		},
	}
}
