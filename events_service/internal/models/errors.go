package models

import (
	"errors"
)

var (
	ErrWrongEventId         = errors.New("no event found")
	ErrWrongUserId          = errors.New("no participant found")
	ErrWrongTimeFormat      = errors.New("wrong time format")
	ErrAlreadyRegistered    = errors.New("user already registered")
	ErrRegistrationNotFound = errors.New("registration not found")
	ErrWrongArgument        = errors.New("wrong argument")
	ErrAccessDenied         = errors.New("access denied")
	ErrGetFromContexxt      = errors.New("failed to get a value from context")
)
