package store

import (
	"context"
	"database/sql"
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
	db *sql.DB
}

func (s *CommentStore) GetAllByPostID(ctx context.Context, postID int64) ([]Comment, error) {
	var query = `
	SELECT
	c.id, c.post_id, c.user_id, c.content, c.created_at, users.username, users.id
	FROM
	comments c
	JOIN users ON users.id = c.user_id
	WHERE c.post_id = $1
	ORDER BY c.created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment

	for rows.Next() {
		var comment Comment
		comment.User = User{}

		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.Content,
			&comment.CreatedAt,
			&comment.User.Username,
			&comment.User.ID,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}
