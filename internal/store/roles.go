package store

import (
	"context"
	"github.com/sergdort/Social/business/domain"
	"github.com/sergdort/Social/internal/store/sqlc"
)

type RolesStore struct {
	queries *sqlc.Queries
}

func (s *RolesStore) GetByRoleType(ctx context.Context, name domain.RoleType) (*domain.Role, error) {
	row, err := s.queries.GetRoleByName(ctx, string(name))
	if err != nil {
		return nil, err
	}
	return &domain.Role{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description.String,
		Level:       int64(row.Level),
	}, nil
}
