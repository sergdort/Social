package domain

import (
	"errors"
)

var ErrNotFound = errors.New("record not found")
var ErrDuplicateEmail = errors.New("email already exists")
var ErrDuplicateUsername = errors.New("username already exists")
