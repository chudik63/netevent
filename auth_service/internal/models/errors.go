package models

import "errors"

var (
	ErrUserAlreadyExists = errors.New("User with such name already exists")
	ErrUserNotFound      = errors.New("User does not exist")
	ErrSignatureInvalid  = errors.New("signature is invalid")
	ErrTokenExpired      = errors.New("token is expired")
	ErrGetFromClaims     = errors.New("argument not found in claims or invalid type")
)
