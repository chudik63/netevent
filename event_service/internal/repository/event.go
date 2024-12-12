package repository

import (
	"context"
	"strconv"

	"gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/database/postgres"
	"gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
)

type Creds map[string]interface{}

type EventRepository struct {
	db postgres.DB
}

func New(db postgres.DB) *EventRepository {
	return &EventRepository{db}
}

func (r *EventRepository) CreateEvent(ctx context.Context, event *models.Event) (int64, error) {
	var id int64

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	err = sq.Insert("public.events").
		Columns("creator_id", "title", "description", "time", "place").
		Values(event.CreatorID, event.Title, event.Description, event.Time, event.Place).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		RunWith(tx).
		QueryRow().
		Scan(&id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if len(event.Topics) > 0 {
		insert := sq.Insert("public.topics").
			Columns("event_id", "topic").
			PlaceholderFormat(sq.Dollar).
			RunWith(tx)

		for _, topic := range event.Topics {
			insert = insert.Values(id, topic)
		}

		_, err = insert.Exec()
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *EventRepository) ReadEvent(ctx context.Context, eventID int64) (*models.Event, error) {
	var event models.Event

	var topics pq.StringArray

	err := sq.Select(
		"e.id AS event_id",
		"e.creator_id",
		"e.title",
		"e.description",
		"e.time",
		"e.place",
		"COALESCE(array_agg(t.topic), '{}') AS topics",
	).
		From("public.events e").
		LeftJoin("public.topics t ON e.id = t.event_id").
		Where(sq.Eq{"e.id": strconv.FormatInt(eventID, 10)}).
		GroupBy("e.id, e.creator_id, e.title, e.description, e.time, e.place").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		QueryRow().
		Scan(&event.EventID, &event.CreatorID, &event.Title, &event.Description, &event.Time, &event.Place, &topics)

	event.Topics = []string(topics)

	return &event, err
}

func (r *EventRepository) UpdateEvent(ctx context.Context, event *models.Event) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = sq.Update("public.events").
		Set("creator_id", event.CreatorID).
		Set("title", event.Title).
		Set("description", event.Description).
		Set("time", event.Time).
		Set("place", event.Place).
		Where(sq.Eq{"id": event.EventID}).
		PlaceholderFormat(sq.Dollar).
		RunWith(tx).
		Exec()
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = sq.Delete("public.topics").
		Where(sq.Eq{"event_id": event.EventID}).
		PlaceholderFormat(sq.Dollar).
		RunWith(tx).
		Exec()
	if err != nil {
		tx.Rollback()
		return err
	}

	if len(event.Topics) > 0 {
		insert := sq.Insert("public.topics").
			Columns("event_id", "topic").
			PlaceholderFormat(sq.Dollar).
			RunWith(tx)

		for _, topic := range event.Topics {
			insert = insert.Values(event.EventID, topic)
		}

		_, err = insert.Exec()
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *EventRepository) DeleteEvent(ctx context.Context, eventID int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = sq.Delete("public.topics").
		Where(sq.Eq{"event_id": eventID}).
		PlaceholderFormat(sq.Dollar).
		RunWith(tx).
		Exec()
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = sq.Delete("public.events").
		Where(sq.Eq{"id": eventID}).
		PlaceholderFormat(sq.Dollar).
		RunWith(tx).
		Exec()
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *EventRepository) ListEvents(ctx context.Context, equations Creds) ([]*models.Event, error) {
	query := sq.Select("e.id, e.creator_id, e.title, e.description, e.time, e.place, COALESCE(array_agg(t.topic), '{}') AS topics").
		From("public.events e").
		LeftJoin("public.topics t ON e.id = t.event_id").
		GroupBy("e.id").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db)

	for key, value := range equations {
		query = query.Where(sq.Eq{key: value})
	}

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*models.Event

	for rows.Next() {
		var event models.Event
		var topics pq.StringArray

		if err := rows.Scan(&event.EventID, &event.CreatorID, &event.Title, &event.Description, &event.Time, &event.Place, &topics); err != nil {
			return nil, err
		}

		event.Topics = []string(topics)

		events = append(events, &event)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return events, nil
}

func (r *EventRepository) RegisterUser(ctx context.Context, userID int64, eventID int64) error {
	_, err := sq.Insert("public.registrations").
		Columns("event_id", "participant_id").
		Values(eventID, userID).
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
	rows, err := sq.Select("public.participants.user_id, public.participants.name").
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

func (r *EventRepository) ListEventsByInterests(ctx context.Context, userID int64) ([]*models.Event, error) {

	return []*models.Event{}, nil
}
