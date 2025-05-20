package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/sergdort/Social/business/domain"
	sqlc2 "github.com/sergdort/Social/business/platform/store/sqlc"
	"time"
)

type UserStore struct {
	db      *sql.DB
	queries *sqlc2.Queries
}

func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)

	defer cancel()
	result, err := s.queries.WithTx(tx).CreateUser(
		ctx,
		sqlc2.CreateUserParams{
			Username: user.Username,
			Email:    user.Email,
			Password: user.Password.Hash,
			RoleID:   int32(user.RoleID),
		},
	)

	user.ID = result.ID
	user.CreatedAt = result.CreatedAt.String()

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return domain.ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return domain.ErrDuplicateUsername
		default:
			return err
		}
	}
	return nil
}

func (s *UserStore) RevertCreateAndInvite(ctx context.Context, id int64) error {
	return withTransaction(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.deleteUser(ctx, id, tx); err != nil {
			return err
		}
		if err := s.deleteInvitation(ctx, id, tx); err != nil {
			return err
		}
		return nil
	})
}

func (s *UserStore) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	row, err := s.queries.GetUserByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, domain.ErrNotFound
		default:
			return nil, err
		}
	}

	return &domain.User{
		ID:        row.ID,
		Username:  row.Username,
		Email:     row.Email,
		CreatedAt: row.CreatedAt.String(),
		IsActive:  row.IsActive,
		RoleID:    row.RoleID,
		Role: domain.Role{
			ID:          row.RoleID,
			Name:        row.RoleName,
			Description: row.RoleDescription.String,
			Level:       int64(row.RoleLevel),
		},
	}, nil
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	row, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	user := domain.User{
		ID:        row.ID,
		Username:  row.Username,
		Email:     row.Email,
		CreatedAt: row.CreatedAt.String(),
		IsActive:  row.IsActive,
		Password: domain.Password{
			Hash: row.Password,
		},
	}
	return &user, nil
}

func (s *UserStore) CreateAndInvite(ctx context.Context, user *domain.User, token string, expiration time.Duration) error {
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

func (s *UserStore) Activate(ctx context.Context, token string) error {
	return withTransaction(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.activateUserByInvitationToken(ctx, token, tx); err != nil {
			return err
		}
		if err := s.deleteUserInvitation(ctx, tx, token); err != nil {
			return err
		}
		return nil
	})
}

func (s *UserStore) activateUserByInvitationToken(ctx context.Context, token string, tx *sql.Tx) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.queries.WithTx(tx).ActiveUserByInvitationToken(ctx, sqlc2.ActiveUserByInvitationTokenParams{
		Token:  []byte(token),
		Expiry: time.Now(),
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return domain.ErrNotFound
		default:
			return err
		}
	}

	return nil
}

func (s *UserStore) createUserInvitation(
	ctx context.Context,
	tx *sql.Tx,
	token string,
	userID int64,
	expiration time.Duration,
) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.queries.WithTx(tx).CreateUserInvitation(ctx, sqlc2.CreateUserInvitationParams{
		Token:  []byte(token),
		UserID: userID,
		Expiry: time.Now().Add(expiration),
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) deleteUserInvitation(ctx context.Context, tx *sql.Tx, token string) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	err := s.queries.WithTx(tx).DeleteUserInvitationByToken(ctx, []byte(token))
	return err
}

func (s *UserStore) deleteUser(ctx context.Context, id int64, tx *sql.Tx) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	err := s.queries.WithTx(tx).DeleteUserByID(ctx, id)
	return err
}

func (s *UserStore) deleteInvitation(ctx context.Context, userID int64, tx *sql.Tx) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	err := s.queries.WithTx(tx).DeleteUserInvitationByUserID(ctx, userID)
	return err
}
