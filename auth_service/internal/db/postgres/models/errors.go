package models

import "errors"

var (
	ErrUserAlreadyExists = errors.New("User with such name already exists")
	ErrUserNotFound      = errors.New("User does not exist")
)
