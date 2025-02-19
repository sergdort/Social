package store

import (
	"context"
	"database/sql"
)

type User struct {
	ID        int64  `json:"id"`
	Username  string `db:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	var query = `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, created_at`

	return s.db.QueryRowContext(
		ctx, query, user.Username, user.Email, user.Password,
	).Scan(&user.ID, &user.CreatedAt)
}
