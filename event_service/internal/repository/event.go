package repository

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/chudik63/netevent/event_service/internal/database/postgres"
	"github.com/chudik63/netevent/event_service/internal/models"

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

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrWrongEventId
		}
		return nil, err
	}

	event.Topics = []string(topics)

	return &event, nil
}

func (r *EventRepository) UpdateEvent(ctx context.Context, event *models.Event) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	res, err := sq.Update("public.events").
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

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if rowsAffected == 0 {
		tx.Rollback()
		return models.ErrWrongEventId
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

	res, err := sq.Delete("public.events").
		Where(sq.Eq{"id": eventID}).
		PlaceholderFormat(sq.Dollar).
		RunWith(tx).
		Exec()
	if err != nil {
		tx.Rollback()
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if rowsAffected == 0 {
		tx.Rollback()
		return models.ErrWrongEventId
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
		GroupBy("e.id, e.creator_id, e.title, e.description, e.time, e.place").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db)

	for key, value := range equations {
		if key != "id" && key != "creator_id" && key != "title" && key != "description" && key != "time" && key != "place" {
			return nil, models.ErrWrongArgument
		}
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

func (r *EventRepository) CreateRegistration(ctx context.Context, userID int64, eventID int64) error {
	_, err := sq.Insert("public.registrations").
		Columns("event_id", "user_id").
		Values(eventID, userID).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		Exec()

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {

			return models.ErrAlreadyRegistered
		}

		return err
	}

	return nil
}

func (r *EventRepository) ListRegistratedEvents(ctx context.Context, userID int64) ([]*models.Event, error) {
	rows, err := sq.Select("e.id, e.creator_id, e.title, e.description, e.time, e.place, COALESCE(array_agg(t.topic), '{}') AS topics").
		From("public.registrations r").
		Join("public.events e ON r.event_id = e.id").
		LeftJoin("public.topics t ON e.id = t.event_id").
		Where(sq.Eq{"r.user_id": userID}).
		GroupBy("e.id, e.creator_id, e.title, e.description, e.time, e.place").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		Query()
	if err != nil {
		return nil, err
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

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
		return nil, rows.Err()
	}

	return events, nil
}

func (r *EventRepository) CreateParticipant(ctx context.Context, participant *models.Participant) error {
	var id int64

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = sq.Insert("public.users").
		Columns("user_id", "name", "email").
		Values(participant.UserID, participant.Name, participant.Email).
		Suffix("RETURNING user_id").
		PlaceholderFormat(sq.Dollar).
		RunWith(tx).
		QueryRow().
		Scan(&id)
	if err != nil {
		tx.Rollback()
		return err
	}

	if len(participant.Interests) > 0 {
		insert := sq.Insert("public.interests").
			Columns("user_id", "interest").
			PlaceholderFormat(sq.Dollar).
			RunWith(tx)

		for _, interest := range participant.Interests {
			insert = insert.Values(id, interest)
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

func (r *EventRepository) ReadParticipant(ctx context.Context, userID int64) (*models.Participant, error) {
	var user models.Participant
	var interests pq.StringArray

	err := sq.Select(
		"p.user_id AS user_id",
		"p.name",
		"p.email",
		"COALESCE(array_agg(i.interest), '{}') AS interests",
	).
		From("public.users p").
		LeftJoin("public.interests i ON p.user_id = i.user_id").
		Where(sq.Eq{"p.user_id": strconv.FormatInt(userID, 10)}).
		GroupBy("p.user_id, p.name, p.email").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		QueryRow().
		Scan(&user.UserID, &user.Name, &user.Email, &interests)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrWrongUserId
		}

		return nil, err
	}
	user.Interests = []string(interests)

	return &user, nil
}

func (r *EventRepository) UpdateParticipant(ctx context.Context, participant *models.Participant) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	res, err := sq.Update("public.users").
		Set("name", participant.Name).
		Set("email", participant.Email).
		Where(sq.Eq{"user_id": participant.UserID}).
		PlaceholderFormat(sq.Dollar).
		RunWith(tx).
		Exec()
	if err != nil {
		tx.Rollback()
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if rowsAffected == 0 {
		tx.Rollback()
		return models.ErrWrongUserId
	}

	_, err = sq.Delete("public.interests").
		Where(sq.Eq{"user_id": participant.UserID}).
		PlaceholderFormat(sq.Dollar).
		RunWith(tx).
		Exec()
	if err != nil {
		tx.Rollback()
		return err
	}

	if len(participant.Interests) > 0 {
		insert := sq.Insert("public.interests").
			Columns("user_id", "interest").
			PlaceholderFormat(sq.Dollar).
			RunWith(tx)

		for _, interest := range participant.Interests {
			insert = insert.Values(participant.UserID, interest)
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

func (r *EventRepository) SetChatStatus(ctx context.Context, userID int64, eventID int64, isReady bool) error {
	_, err := sq.Update("public.registrations").
		Set("ready_to_chat", isReady).
		Where(sq.Eq{"event_id": strconv.FormatInt(eventID, 10), "user_id": strconv.FormatInt(userID, 10)}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		Exec()

	return err
}

func (r *EventRepository) ListUsersToChat(ctx context.Context, eventID int64, userID int64) ([]*models.Participant, error) {
	rows, err := sq.Select("r.user_id, u.name, COALESCE(array_agg(i.interest), '{}')").
		From("public.registrations r").
		LeftJoin("public.users u ON r.user_id = u.user_id").
		LeftJoin("public.interests i ON r.user_id = i.user_id").
		Where(sq.Eq{"r.event_id": strconv.FormatInt(eventID, 10), "r.ready_to_chat": true}).
		Where("u.user_id != ?", userID).
		GroupBy("r.user_id, u.name").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		Query()

	if err != nil {
		return nil, err
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	var users []*models.Participant
	for rows.Next() {
		var user models.Participant
		if err := rows.Scan(&user.UserID, &user.Name, pq.Array(&user.Interests)); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return users, nil
}

func (r *EventRepository) ListEventsByInterests(ctx context.Context, userID int64) ([]*models.Event, error) {
	var interests []string
	err := sq.Select("COALESCE(array_agg(i.interest), '{}')").
		From("public.interests i").
		Join("public.users p ON p.user_id = i.user_id").
		Where(sq.Eq{"i.user_id": userID}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		QueryRow().
		Scan(pq.Array(&interests))

	if err != nil {
		return nil, err
	}

	var events []*models.Event
	query := sq.Select("e.id, e.creator_id, e.title, e.description, e.time, e.place, COALESCE(array_agg(t.topic), '{}') AS topics").
		From("public.events e").
		LeftJoin("public.topics t ON e.id = t.event_id").
		PlaceholderFormat(sq.Dollar).
		GroupBy("e.id, e.creator_id, e.title, e.description, e.time, e.place").
		RunWith(r.db)

	for _, interest := range interests {
		query = query.Where(sq.Like{"t.topic": "%" + interest + "%"})
	}

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

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
