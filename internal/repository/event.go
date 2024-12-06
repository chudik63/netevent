package repository

import (
	"context"
	"event_service/internal/database/postgres"
	"event_service/internal/models"
	"strconv"

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
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		QueryRow().
		Scan(&id)

	return id, err
}

func (r *EventRepository) ReadEvent(ctx context.Context, eventID int64) (*models.Event, error) {
	var event models.Event

	err := sq.Select("*").
		From("public.events").
		Where(sq.Eq{"id": strconv.FormatInt(eventID, 10)}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		QueryRow().
		Scan(&event.EventID, &event.CreatorID, &event.Title, &event.Description, &event.Time, &event.Place)

	return &event, err
}

func (r *EventRepository) UpdateEvent(ctx context.Context, event *models.Event) error {
	_, err := sq.Update("public.events").
		Set("creator_id", event.CreatorID).
		Set("title", event.Title).
		Set("description", event.Description).
		Set("time", event.Time).
		Set("place", event.Place).
		Where(sq.Eq{"id": event.EventID}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		Exec()

	return err
}

func (r *EventRepository) DeleteEvent(ctx context.Context, eventID int64) error {
	_, err := sq.Delete("public.events").
		Where(sq.Eq{"id": strconv.FormatInt(eventID, 10)}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		Exec()

	return err
}

func (r *EventRepository) ListEventsByCreator(ctx context.Context, creatorID int64) ([]*models.Event, error) {
	rows, err := sq.Select("*").
		From("public.events").
		Where(sq.Eq{"creator_id": strconv.FormatInt(creatorID, 10)}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		Query()
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var events []*models.Event
	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.EventID, &event.CreatorID, &event.Title, &event.Description, &event.Time, &event.Place); err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return events, nil
}
