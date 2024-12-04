package repository

import (
	"context"
	"event_service/internal/database/postgres"
	"event_service/internal/models"

	sq "github.com/Masterminds/squirrel"
)

type EventRepository struct {
	db postgres.DB
}

func New(db postgres.DB) *EventRepository {
	return &EventRepository{db}
}

func (r *EventRepository) CreateEvent(ctx context.Context, event *models.Event) (int64, error) {
	var id int64

	err := sq.Insert("public.events").
		Columns("creator_id", "title", "description", "time", "place").
		Values(event.CreatorID, event.Title, event.Description, event.Time, event.Place).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		QueryRow().
		Scan(&id)

	if err != nil {
		return id, err
	}

	return id, nil
}

func (r *EventRepository) ReadEvent(ctx context.Context, eventID int64) (*models.Event, error) {
	var event models.Event

	err := sq.Select("*").
		From("public.events").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		QueryRow().
		Scan(&event)

	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *EventRepository) UpdateEvent(ctx context.Context, event *models.Event) error {
	return nil
}

func (r *EventRepository) DeleteEvent(ctx context.Context, eventID int64) error {
	return nil
}

func (r *EventRepository) ListEvents(ctx context.Context) []models.Event {
	return []models.Event{}
}
