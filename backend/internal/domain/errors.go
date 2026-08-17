package domain

import (
	"errors"
)

var (
	ErrAlreadyExists = errors.New("short code already exists")
	ErrNotFound      = errors.New("url not found")
)
