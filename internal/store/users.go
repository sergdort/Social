package store

import (
	"context"
	"database/sql"
	"errors"
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

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)

	defer cancel()

	return s.db.QueryRowContext(
		ctx, query, user.Username, user.Email, user.Password,
	).Scan(&user.ID, &user.CreatedAt)
}

func (s *UserStore) GetByID(ctx context.Context, id int64) (*User, error) {
	var query = `SELECT id, username, email, created_at FROM users WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var user User
	err := s.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
