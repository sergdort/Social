package store

import (
	"context"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID        int64    `json:"id"`
	Username  string   `db:"username"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
	CreatedAt string   `json:"created_at"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.hash = hash
	p.text = &text

	return nil
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	var query = `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)

	defer cancel()
	var err error
	if tx == nil {
		err = s.db.QueryRowContext(
			ctx,
			query,
			user.Username,
			user.Email,
			user.Password.hash,
		).Scan(
			&user.ID,
			&user.CreatedAt,
		)
	} else {
		err = tx.QueryRowContext(
			ctx,
			query,
			user.Username,
			user.Email,
			user.Password.hash,
		).Scan(
			&user.ID,
			&user.CreatedAt,
		)
	}

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}
	return nil
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

func (s *UserStore) CreateAndInvite(ctx context.Context, user *User, token string, expiration time.Duration) error {
	return withTransaction(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}

		if err := s.createUserInvitation(ctx, tx, token, user.ID, expiration); err != nil {
			return err
		}
		return nil
	})
}

func (s *UserStore) createUserInvitation(
	ctx context.Context,
	tx *sql.Tx,
	token string,
	userID int64,
	expiration time.Duration,
) error {
	query := `INSERT INTO user_invitations (token, user_id, expiry) VALUES ($1, $2, $3)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	if _, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(expiration)); err != nil {
		return err
	}
	return nil
}
