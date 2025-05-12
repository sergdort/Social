package usersapp

import (
	"encoding/json"
	"github.com/sergdort/Social/business/domain"
)

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	IsActive  bool   `json:"is_active"`
	Role      string `json:"role"`
}

func toAppUser(domain *domain.User) User {
	return User{
		ID:        domain.ID,
		Username:  domain.Username,
		Email:     domain.Email,
		CreatedAt: domain.CreatedAt,
		IsActive:  domain.IsActive,
		Role:      domain.Role.Name,
	}
}

func (usr User) Encode() ([]byte, string, error) {
	data, err := json.Marshal(usr)
	return data, "application/json", err
}
