package domain

import "context"

type Comment struct {
	ID        int64  `json:"id"`
	PostID    int64  `json:"post_id"`
	UserID    int64  `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`

	User User `json:"user"`
}

type CommentsRepository interface {
	Create(ctx context.Context, comment *Comment) error
	GetAllByPostID(ctx context.Context, postID int64) ([]Comment, error)
}
