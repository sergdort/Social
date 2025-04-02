package store

import (
	"context"
	"github.com/sergdort/Social/internal/store/sqlc"
)

type RoleType string

type Role struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Level       int64  `json:"level"`
}

// Allowed values for RoleType
const (
	RoleTypeUser      RoleType = "user"
	RoleTypeModerator RoleType = "moderator"
	RoleTypeAdmin     RoleType = "admin"
)

type RolesStore struct {
	queries *sqlc.Queries
}

func (s *RolesStore) GetByRoleType(ctx context.Context, name RoleType) (*Role, error) {
	row, err := s.queries.GetRoleByName(ctx, string(name))
	if err != nil {
		return nil, err
	}
	return &Role{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description.String,
		Level:       int64(row.Level),
	}, nil
}

// IsValid checks if the role is one of the allowed values
func (r RoleType) IsValid() bool {
	switch r {
	case RoleTypeUser, RoleTypeModerator, RoleTypeAdmin:
		return true
	default:
		return false
	}
}
