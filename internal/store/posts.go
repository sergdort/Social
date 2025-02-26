package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
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
	CommentsCount int `json:"comments_count"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	var query = `INSERT INTO posts (content, title, user_id, tags) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`

	return s.db.QueryRowContext(
		ctx, query, post.Content, post.Title, post.UserID, pq.Array(post.Tags),
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
}

func (s *PostStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	var query = `SELECT id, content, title, user_id, created_at, updated_at, tags, version FROM posts WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)

	defer cancel()

	var post Post
	var err = s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Content,
		&post.Title,
		&post.UserID,
		&post.CreatedAt,
		&post.UpdatedAt,
		pq.Array(&post.Tags),
		&post.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &post, nil
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
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
	SELECT
		p.id,
		p.user_id,
		p.title,
		p.content,
		p.created_at,
		p.tags,
		COUNT(c.id) AS comments_count,
		u.username
	FROM
		posts p
		LEFT JOIN comments c ON c.post_id = p.id
		LEFT JOIN users u ON p.user_id = u.id
		JOIN followers f ON f.follower_id = p.user_id
		OR p.user_id = $1 OR p.user_id = $1
	WHERE
		(f.user_id = $1 OR p.user_id = $1) AND
		(p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%') AND
		(p.tags @> $5 OR $5 = '{}')
	GROUP BY
		p.id, u.username
	ORDER BY
		p.created_at ` + q.SortBy + ` LIMIT $2 OFFSET $3`

	rows, err := s.db.QueryContext(ctx, query, userId, q.Limit, q.Offset, q.Search, pq.Array(q.Tags))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var postsWithMetadata []PostWithMetadata

	for rows.Next() {
		var postWithMetadata PostWithMetadata
		var user User
		err = rows.Scan(
			&postWithMetadata.Post.ID,
			&postWithMetadata.Post.UserID,
			&postWithMetadata.Post.Title,
			&postWithMetadata.Post.Content,
			&postWithMetadata.Post.CreatedAt,
			pq.Array(&postWithMetadata.Post.Tags),
			&postWithMetadata.CommentsCount,
			&user.Username,
		)
		user.ID = postWithMetadata.Post.UserID
		postWithMetadata.Post.User = user

		if err != nil {
			return nil, err
		}
		postsWithMetadata = append(postsWithMetadata, postWithMetadata)
	}

	return postsWithMetadata, nil
}
