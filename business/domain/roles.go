package domain

import "context"

type Role struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Level       int64  `json:"level"`
}

type RoleType string

// Allowed values for RoleType
const (
	RoleTypeUser      RoleType = "user"
	RoleTypeModerator RoleType = "moderator"
	RoleTypeAdmin     RoleType = "admin"
)

type RolesRepository interface {
	GetByRoleType(ctx context.Context, name RoleType) (*Role, error)
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
