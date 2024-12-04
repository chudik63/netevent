package repository

import (
	"event_service/internal/database/postgres"
)

type EventRepository struct {
	db postgres.DB
}

func New(db postgres.DB) *EventRepository {
	return &EventRepository{db}
}
