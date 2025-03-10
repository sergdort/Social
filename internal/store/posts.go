package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"github.com/sergdort/Social/internal/store/sqlc"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Tags      []string  `json:"tags"`
	Comments  []Comment `json:"comments"`
	Version   int64     `json:"version"`
	User      User      `json:"user"`
}

type PostWithMetadata struct {
	Post
	CommentsCount int64 `json:"comments_count"`
}

type PostStore struct {
	db      *sql.DB
	queries *sqlc.Queries
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	var query = `INSERT INTO posts (content, title, user_id, tags) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`

	return s.db.QueryRowContext(
		ctx, query, post.Content, post.Title, post.UserID, pq.Array(post.Tags),
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
}

func (s *PostStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)

	defer cancel()

	row, err := s.queries.GetPostByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &Post{
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
	var query = `DELETE FROM posts WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)

	defer cancel()

	res, err := s.db.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *PostStore) Update(ctx context.Context, post *Post) error {
	query := `
		UPDATE posts 
		SET content = $1, title = $2, version = version + 1
		WHERE id = $3 AND version = $4 
		RETURNING version
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)

	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.ID,
		post.Version,
	).Scan(&post.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}

func (s *PostStore) GetUserFeed(ctx context.Context, userId int64, q PaginatedFeedQuery) ([]PostWithMetadata, error) {
	feed, err := s.queries.GetUserFeed(ctx, sqlc.GetUserFeedParams{
		UserID:  userId,
		Limit:   int32(q.Limit),
		Offset:  int32(q.Offset),
		Column4: q.Search,
		Tags:    q.Tags,
	})
	if err != nil {
		return nil, err
	}
	postsWithMetadata := Map(feed, convertToPostWithMetadata)
	return postsWithMetadata, nil
}

func convertToPostWithMetadata(feedRow sqlc.GetUserFeedRow) PostWithMetadata {
	return PostWithMetadata{
		Post: Post{
			ID:        feedRow.ID,
			Content:   feedRow.Content,
			Title:     feedRow.Title,
			UserID:    feedRow.UserID,
			CreatedAt: feedRow.CreatedAt.String(),
			Tags:      feedRow.Tags,
			User: User{
				ID:       feedRow.UserID,
				Username: feedRow.Username.String,
			},
		},
		CommentsCount: feedRow.CommentsCount,
	}
}

func Map[T any, U any](input []T, mapper func(T) U) []U {
	result := make([]U, len(input))
	for i, v := range input {
		result[i] = mapper(v)
	}
	return result
}
