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

func (r *EventRepository) RegisterUser(ctx context.Context, participant *models.Participant, eventID int64) error {
	_, err := sq.Insert("public.participants").
		Columns("user_id", "name").
		Values(participant.UserID, participant.Name).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		Exec()

	if err != nil {
		return err
	}

	_, err = sq.Insert("public.registrations").
		Columns("event_id", "participant_id").
		Values(eventID, participant.UserID).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		Exec()

	return err
}

func (r *EventRepository) SetChatStatus(ctx context.Context, participantID int64, eventID int64, isReady bool) error {
	_, err := sq.Update("public.registrations").
		Set("ready_to_chat", isReady).
		Where(sq.Eq{"event_id": strconv.FormatInt(eventID, 10), "participant_id": strconv.FormatInt(participantID, 10)}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		Exec()

	return err
}

func (r *EventRepository) ListUsersToChat(ctx context.Context, eventID int64) ([]*models.Participant, error) {
	rows, err := sq.Select("*").
		From("public.registrations").
		LeftJoin("public.participants ON public.registrations.participant_id = public.participants.id").
		Where(sq.Eq{"public.registrations.event_id": strconv.FormatInt(eventID, 10), "public.registrations.ready_to_chat": "true"}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		Query()
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var users []*models.Participant
	for rows.Next() {
		var user models.Participant
		if err := rows.Scan(&user.UserID, &user.Name); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return users, nil
}

func (r *EventRepository) ListEventsByUser(ctx context.Context, userID int64) ([]*models.Event, error) {
	rows, err := sq.Select("public.events.id", "creator_id", "title", "description", "time", "place").
		From("public.registrations").
		LeftJoin("public.events ON public.events.id = public.registrations.event_id").
		LeftJoin("public.participants ON public.participants.id = public.registrations.participant_id").
		Where(sq.Eq{"public.participants.user_id": strconv.FormatInt(userID, 10)}).
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

func (r *EventRepository) ListEventsByInterests(ctx context.Context, userID int64) ([]*models.Event, error) {
	return []*models.Event{}, nil
}
