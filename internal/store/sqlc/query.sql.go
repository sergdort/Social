// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package sqlc

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

const getPostByID = `-- name: GetPostByID :one
SELECT id,
       content,
       title,
       user_id,
       created_at,
       updated_at,
       tags,
       version
FROM posts
WHERE id = $1
`

type GetPostByIDRow struct {
	ID        int64
	Content   string
	Title     string
	UserID    int64
	CreatedAt time.Time
	UpdatedAt time.Time
	Tags      []string
	Version   sql.NullInt32
}

func (q *Queries) GetPostByID(ctx context.Context, id int64) (GetPostByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getPostByID, id)
	var i GetPostByIDRow
	err := row.Scan(
		&i.ID,
		&i.Content,
		&i.Title,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
		pq.Array(&i.Tags),
		&i.Version,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, email, password, username, created_at, is_active
FROM users
WHERE email = $1
`

type GetUserByEmailRow struct {
	ID        int64
	Email     string
	Password  []byte
	Username  string
	CreatedAt time.Time
	IsActive  bool
}

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (GetUserByEmailRow, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i GetUserByEmailRow
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.Username,
		&i.CreatedAt,
		&i.IsActive,
	)
	return i, err
}

const getUserFeed = `-- name: GetUserFeed :many
SELECT p.id,
       p.user_id,
       p.title,
       p.content,
       p.created_at,
       p.tags,
       COUNT(c.id) AS comments_count,
       u.username
FROM posts p
         LEFT JOIN comments c ON c.post_id = p.id
         LEFT JOIN users u ON p.user_id = u.id
         JOIN followers f ON f.follower_id = p.user_id
    OR p.user_id = $1
WHERE (f.user_id = $1 OR p.user_id = $1)
  AND ($4 = '' OR LOWER(p.title) LIKE LOWER('%' || $4 || '%') OR LOWER(p.content) LIKE LOWER('%' || $4 || '%'))
  AND (p.tags @> $5 OR $5 = '{}')
GROUP BY p.id, u.username
ORDER BY p.created_at DESC
LIMIT $2 OFFSET $3
`

type GetUserFeedParams struct {
	UserID  int64
	Limit   int32
	Offset  int32
	Column4 interface{}
	Tags    []string
}

type GetUserFeedRow struct {
	ID            int64
	UserID        int64
	Title         string
	Content       string
	CreatedAt     time.Time
	Tags          []string
	CommentsCount int64
	Username      sql.NullString
}

func (q *Queries) GetUserFeed(ctx context.Context, arg GetUserFeedParams) ([]GetUserFeedRow, error) {
	rows, err := q.db.QueryContext(ctx, getUserFeed,
		arg.UserID,
		arg.Limit,
		arg.Offset,
		arg.Column4,
		pq.Array(arg.Tags),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserFeedRow
	for rows.Next() {
		var i GetUserFeedRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Title,
			&i.Content,
			&i.CreatedAt,
			pq.Array(&i.Tags),
			&i.CommentsCount,
			&i.Username,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
