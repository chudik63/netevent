package models

import (
	"errors"
)

var (
	ErrWrongEventId = errors.New("no event found")
	ErrWrongUserId  = errors.New("no participant found")
)
