package repository

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"testing"

	"github.com/chudik63/netevent/event_service/internal/database/postgres"
	"github.com/chudik63/netevent/event_service/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCreateEvent(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	r := New(postgres.DB{DB: db})

	type mockBehavior func(event *models.Event, id int64)

	testTable := []struct {
		name          string
		inputEvent    *models.Event
		mockBehavior  mockBehavior
		expectedId    int64
		expectedError error
	}{
		{
			name: "OK test",
			inputEvent: &models.Event{
				CreatorID:   1,
				Title:       "test event",
				Description: "test description",
				Time:        "2022-01-01 13:31:00",
				Place:       "test place",
				Topics:      []string{"test1", "test2"},
			},
			mockBehavior: func(event *models.Event, id int64) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO public.events").
					WithArgs(event.CreatorID, event.Title, event.Description, event.Time, event.Place).
					WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO public.topics").
					WithArgs(id, "test1", id, "test2").
					WillReturnResult(sqlmock.NewResult(2, 2))

				mock.ExpectCommit()
			},
			expectedId:    1,
			expectedError: nil,
		},
		{
			name:       "Insert Event Error",
			inputEvent: &models.Event{},
			mockBehavior: func(event *models.Event, id int64) {
				mock.ExpectBegin()

				mock.ExpectQuery("INSERT INTO public.events").
					WithArgs(event.CreatorID, event.Title, event.Description, event.Time, event.Place).
					WillReturnError(errors.New("insert event error"))

				mock.ExpectRollback()
			},
			expectedId:    0,
			expectedError: errors.New("insert event error"),
		},
		{
			name: "Insert Topics Error",
			inputEvent: &models.Event{
				CreatorID:   1,
				Title:       "test event",
				Description: "test description",
				Time:        "2022-01-01 13:31:00",
				Place:       "test place",
				Topics:      []string{"test1", "test2"},
			},
			mockBehavior: func(event *models.Event, id int64) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO public.events").
					WithArgs(event.CreatorID, event.Title, event.Description, event.Time, event.Place).
					WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO public.topics").
					WithArgs(id, "test1", id, "test2").
					WillReturnError(errors.New("insert topics error"))

				mock.ExpectRollback()
			},
			expectedId:    0,
			expectedError: errors.New("insert topics error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.inputEvent, testCase.expectedId)

			id, err := r.CreateEvent(context.Background(), testCase.inputEvent)

			if testCase.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedId, id)
			}
		})
	}
}

func TestReadEvent(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	r := New(postgres.DB{DB: db})

	type mockBehavior func(eventID int64)

	testTable := []struct {
		name          string
		eventID       int64
		mockBehavior  mockBehavior
		expectedEvent *models.Event
		expectedError error
	}{
		{
			name:    "OK test",
			eventID: 1,
			mockBehavior: func(eventID int64) {
				rows := sqlmock.NewRows([]string{"event_id", "creator_id", "title", "description", "time", "place", "topics"}).
					AddRow(1, 2, "Test Event", "Test Description", "2024-01-01 10:00:00", "Test Place", pq.StringArray{"Topic1", "Topic2"})
				mock.ExpectQuery("SELECT e.id AS event_id, e.creator_id, e.title, e.description, e.time, e.place, .*").
					WithArgs(strconv.FormatInt(eventID, 10)).
					WillReturnRows(rows)
			},
			expectedEvent: &models.Event{
				EventID:     1,
				CreatorID:   2,
				Title:       "Test Event",
				Description: "Test Description",
				Time:        "2024-01-01 10:00:00",
				Place:       "Test Place",
				Topics:      []string{"Topic1", "Topic2"},
			},
			expectedError: nil,
		},
		{
			name:    "Event Not Found",
			eventID: 1,
			mockBehavior: func(eventID int64) {
				mock.ExpectQuery("SELECT e.id AS event_id, e.creator_id, e.title, e.description, e.time, e.place, .*").
					WithArgs(strconv.FormatInt(eventID, 10)).
					WillReturnError(sql.ErrNoRows)
			},
			expectedEvent: nil,
			expectedError: models.ErrWrongEventId,
		},
		{
			name:    "Query Error",
			eventID: 1,
			mockBehavior: func(eventID int64) {
				mock.ExpectQuery("SELECT e.id AS event_id, e.creator_id, e.title, e.description, e.time, e.place, .*").
					WithArgs(strconv.FormatInt(eventID, 10)).
					WillReturnError(errors.New("query error"))
			},
			expectedEvent: nil,
			expectedError: errors.New("query error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.eventID)

			event, err := r.ReadEvent(context.Background(), testCase.eventID)

			if testCase.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedEvent, event)
			}
		})
	}
}

func TestUpdateEvent(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	r := New(postgres.DB{DB: db})

	type mockBehavior func(event *models.Event)

	testTable := []struct {
		name          string
		inputEvent    *models.Event
		mockBehavior  mockBehavior
		expectedError error
	}{
		{
			name: "OK test",
			inputEvent: &models.Event{
				EventID:     1,
				CreatorID:   2,
				Title:       "Updated Title",
				Description: "Updated Description",
				Time:        "2024-01-02 10:00:00",
				Place:       "Updated Place",
				Topics:      []string{"Topic3", "Topic4"},
			},
			mockBehavior: func(event *models.Event) {
				mock.ExpectBegin()

				mock.ExpectExec("UPDATE public.events").
					WithArgs(event.CreatorID, event.Title, event.Description, event.Time, event.Place, event.EventID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("DELETE FROM public.topics").
					WithArgs(event.EventID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("INSERT INTO public.topics").
					WithArgs(event.EventID, "Topic3", event.EventID, "Topic4").
					WillReturnResult(sqlmock.NewResult(2, 2))

				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name:       "Update Event Error",
			inputEvent: &models.Event{EventID: 1},
			mockBehavior: func(event *models.Event) {
				mock.ExpectBegin()

				mock.ExpectExec("UPDATE public.events").
					WithArgs(event.CreatorID, event.Title, event.Description, event.Time, event.Place, event.EventID).
					WillReturnError(errors.New("update event error"))

				mock.ExpectRollback()
			},
			expectedError: errors.New("update event error"),
		},
		{
			name: "Topics Deletion Error",
			inputEvent: &models.Event{
				EventID: 1,
				Topics:  []string{"Topic3"},
			},
			mockBehavior: func(event *models.Event) {
				mock.ExpectBegin()

				mock.ExpectExec("UPDATE public.events").
					WithArgs(event.CreatorID, event.Title, event.Description, event.Time, event.Place, event.EventID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("DELETE FROM public.topics").
					WithArgs(event.EventID).
					WillReturnError(errors.New("delete topics error"))

				mock.ExpectRollback()
			},
			expectedError: errors.New("delete topics error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.inputEvent)

			err := r.UpdateEvent(context.Background(), testCase.inputEvent)

			if testCase.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteEvent(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	r := New(postgres.DB{DB: db})

	type mockBehavior func(eventID int64)

	testTable := []struct {
		name          string
		eventID       int64
		mockBehavior  mockBehavior
		expectedError error
	}{
		{
			name:    "OK test",
			eventID: 1,
			mockBehavior: func(eventID int64) {
				mock.ExpectBegin()

				mock.ExpectExec("DELETE FROM public.topics").
					WithArgs(eventID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("DELETE FROM public.events").
					WithArgs(eventID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name:    "Delete Topics Error",
			eventID: 1,
			mockBehavior: func(eventID int64) {
				mock.ExpectBegin()

				mock.ExpectExec("DELETE FROM public.topics").
					WithArgs(eventID).
					WillReturnError(errors.New("delete topics error"))

				mock.ExpectRollback()
			},
			expectedError: errors.New("delete topics error"),
		},
		{
			name:    "Delete Event Error",
			eventID: 1,
			mockBehavior: func(eventID int64) {
				mock.ExpectBegin()

				mock.ExpectExec("DELETE FROM public.topics").
					WithArgs(eventID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("DELETE FROM public.events").
					WithArgs(eventID).
					WillReturnError(errors.New("delete event error"))

				mock.ExpectRollback()
			},
			expectedError: errors.New("delete event error"),
		},
		{
			name:    "No Rows Affected",
			eventID: 999,
			mockBehavior: func(eventID int64) {
				mock.ExpectBegin()

				mock.ExpectExec("DELETE FROM public.topics").
					WithArgs(eventID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("DELETE FROM public.events").
					WithArgs(eventID).
					WillReturnResult(sqlmock.NewResult(0, 0))

				mock.ExpectRollback()
			},
			expectedError: models.ErrWrongEventId,
		},
		{
			name:    "Transaction Commit Error",
			eventID: 1,
			mockBehavior: func(eventID int64) {
				mock.ExpectBegin()

				mock.ExpectExec("DELETE FROM public.topics").
					WithArgs(eventID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("DELETE FROM public.events").
					WithArgs(eventID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit().WillReturnError(errors.New("commit error"))
			},
			expectedError: errors.New("commit error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.eventID)

			err := r.DeleteEvent(context.Background(), testCase.eventID)

			if testCase.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestListEvents(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	r := New(postgres.DB{DB: db})

	type mockBehavior func(equations Creds)

	testTable := []struct {
		name           string
		equations      Creds
		mockBehavior   mockBehavior
		expectedEvents []*models.Event
		expectedError  error
	}{
		{
			name:      "OK test",
			equations: Creds{"creator_id": int64(1)},
			mockBehavior: func(equations Creds) {
				rows := sqlmock.NewRows([]string{
					"id", "creator_id", "title", "description", "time", "place", "topics",
				}).AddRow(
					1, 1, "Event 1", "Description 1", "2024-12-16 10:00:00", "Place 1", pq.StringArray{"topic1", "topic2"},
				).AddRow(
					2, 1, "Event 2", "Description 2", "2024-12-16 11:00:00", "Place 2", pq.StringArray{"topic3"},
				)

				mock.ExpectQuery("SELECT e.id, e.creator_id, e.title, e.description, e.time, e.place, COALESCE").
					WithArgs(equations["creator_id"]).
					WillReturnRows(rows)
			},
			expectedEvents: []*models.Event{
				{
					EventID:     1,
					CreatorID:   1,
					Title:       "Event 1",
					Description: "Description 1",
					Time:        "2024-12-16 10:00:00",
					Place:       "Place 1",
					Topics:      []string{"topic1", "topic2"},
				},
				{
					EventID:     2,
					CreatorID:   1,
					Title:       "Event 2",
					Description: "Description 2",
					Time:        "2024-12-16 11:00:00",
					Place:       "Place 2",
					Topics:      []string{"topic3"},
				},
			},
			expectedError: nil,
		},
		{
			name:      "Query Error",
			equations: Creds{"creator_id": int64(1)},
			mockBehavior: func(equations Creds) {
				mock.ExpectQuery("SELECT e.id, e.creator_id, e.title, e.description, e.time, e.place, COALESCE").
					WithArgs(equations["creator_id"]).
					WillReturnError(errors.New("query error"))
			},
			expectedEvents: nil,
			expectedError:  errors.New("query error"),
		},
		{
			name:      "Scan Error",
			equations: Creds{"creator_id": int64(1)},
			mockBehavior: func(equations Creds) {
				rows := sqlmock.NewRows([]string{
					"id", "creator_id", "title", "description", "time", "place", "topics",
				}).AddRow(
					nil, 1, "Event 1", "Description 1", "2024-12-16 10:00:00", "Place 1", pq.StringArray{"topic1"},
				)

				mock.ExpectQuery("SELECT e.id, e.creator_id, e.title, e.description, e.time, e.place, COALESCE").
					WithArgs(equations["creator_id"]).
					WillReturnRows(rows)
			},
			expectedEvents: nil,
			expectedError:  errors.New("sql: Scan error on column index 0, name \"id\": converting NULL to int64 is unsupported"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.equations)

			events, err := r.ListEvents(context.Background(), testCase.equations)

			if testCase.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedEvents, events)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Expectations were not met: %v", err)
			}
		})
	}
}

func TestCreateRegistration(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	r := New(postgres.DB{DB: db})

	type mockBehavior func(userID int64, eventID int64)

	testTable := []struct {
		name         string
		userID       int64
		eventID      int64
		mockBehavior mockBehavior
		expectedErr  error
	}{
		{
			name:    "OK test",
			userID:  1,
			eventID: 1,
			mockBehavior: func(userID int64, eventID int64) {
				mock.ExpectExec("INSERT INTO public.registrations").
					WithArgs(eventID, userID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
		{
			name:    "Insert Error",
			userID:  1,
			eventID: 1,
			mockBehavior: func(userID int64, eventID int64) {
				mock.ExpectExec("INSERT INTO public.registrations").
					WithArgs(eventID, userID).
					WillReturnError(errors.New("insert error"))
			},
			expectedErr: errors.New("insert error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.userID, testCase.eventID)

			err := r.CreateRegistration(context.Background(), testCase.userID, testCase.eventID)

			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Expectations were not met: %v", err)
			}
		})
	}
}

func TestCreateParticipant(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	r := New(postgres.DB{DB: db})

	type mockBehavior func(participant *models.Participant)

	testTable := []struct {
		name          string
		participant   *models.Participant
		mockBehavior  mockBehavior
		expectedError error
	}{
		{
			name: "OK test",
			participant: &models.Participant{
				UserID:    1,
				Name:      "Abra Kadabra",
				Email:     "ababa@example.com",
				Interests: []string{"coding", "gaming"},
			},
			mockBehavior: func(participant *models.Participant) {
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO public.users").
					WithArgs(participant.UserID, participant.Name, participant.Email).
					WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))

				mock.ExpectExec("INSERT INTO public.interests").
					WithArgs(1, "coding", 1, "gaming").
					WillReturnResult(sqlmock.NewResult(2, 2))

				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "Users Insert Error",
			participant: &models.Participant{
				UserID:    1,
				Name:      "Abra Kadabra",
				Email:     "ababa@example.com",
				Interests: []string{"coding", "gaming"},
			},
			mockBehavior: func(participant *models.Participant) {
				mock.ExpectBegin()

				mock.ExpectQuery("INSERT INTO public.users").
					WithArgs(participant.UserID, participant.Name, participant.Email).
					WillReturnError(errors.New("insert user error"))

				mock.ExpectRollback()
			},
			expectedError: errors.New("insert user error"),
		},
		{
			name: "Interests Insert Error",
			participant: &models.Participant{
				UserID:    1,
				Name:      "Abra Kadabra",
				Email:     "ababa@example.com",
				Interests: []string{"coding", "gaming"},
			},
			mockBehavior: func(participant *models.Participant) {
				mock.ExpectBegin()

				mock.ExpectQuery("INSERT INTO public.users").
					WithArgs(participant.UserID, participant.Name, participant.Email).
					WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))

				mock.ExpectExec("INSERT INTO public.interests").
					WithArgs(1, "coding", 1, "gaming").
					WillReturnError(errors.New("insert interests error"))

				mock.ExpectRollback()
			},
			expectedError: errors.New("insert interests error"),
		},
		{
			name: "Transaction Commit Error",
			participant: &models.Participant{
				UserID:    1,
				Name:      "Abra Kadabra",
				Email:     "ababa@example.com",
				Interests: []string{"coding", "gaming"},
			},
			mockBehavior: func(participant *models.Participant) {
				mock.ExpectBegin()

				mock.ExpectQuery("INSERT INTO public.users").
					WithArgs(participant.UserID, participant.Name, participant.Email).
					WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))

				mock.ExpectExec("INSERT INTO public.interests").
					WithArgs(1, "coding", 1, "gaming").
					WillReturnResult(sqlmock.NewResult(2, 2))

				mock.ExpectCommit().WillReturnError(errors.New("commit error"))
			},
			expectedError: errors.New("commit error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.participant)

			err := r.CreateParticipant(context.Background(), testCase.participant)

			if testCase.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Expectations were not met: %v", err)
			}
		})
	}
}

func TestReadParticipant(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	r := New(postgres.DB{DB: db})

	type mockBehavior func(userID int64)

	testTable := []struct {
		name                string
		userID              int64
		mockBehavior        mockBehavior
		expectedParticipant *models.Participant
		expectedError       error
	}{
		{
			name:   "OK test",
			userID: 1,
			mockBehavior: func(userID int64) {
				rows := sqlmock.NewRows([]string{
					"user_id", "name", "email", "interests",
				}).AddRow(
					1, "User 1", "user1@example.com", pq.StringArray{"interest1", "interest2"},
				)

				mock.ExpectQuery("SELECT p.user_id AS user_id, p.name, p.email, COALESCE").
					WithArgs(strconv.FormatInt(userID, 10)).
					WillReturnRows(rows)
			},
			expectedParticipant: &models.Participant{
				UserID:    1,
				Name:      "User 1",
				Email:     "user1@example.com",
				Interests: []string{"interest1", "interest2"},
			},
			expectedError: nil,
		},
		{
			name:   "User not found",
			userID: 1,
			mockBehavior: func(userID int64) {
				mock.ExpectQuery("SELECT p.user_id AS user_id, p.name, p.email, COALESCE").
					WithArgs(strconv.FormatInt(userID, 10)).
					WillReturnError(sql.ErrNoRows)
			},
			expectedParticipant: nil,
			expectedError:       models.ErrWrongUserId,
		},
		{
			name:   "Query Error",
			userID: 1,
			mockBehavior: func(userID int64) {
				mock.ExpectQuery("SELECT p.user_id AS user_id, p.name, p.email, COALESCE").
					WithArgs(strconv.FormatInt(userID, 10)).
					WillReturnError(errors.New("query error"))
			},
			expectedParticipant: nil,
			expectedError:       errors.New("query error"),
		},
		{
			name:   "Scan Error",
			userID: 1,
			mockBehavior: func(userID int64) {
				rows := sqlmock.NewRows([]string{
					"user_id", "name", "email", "interests",
				}).AddRow(
					nil, "User 1", "user1@example.com", pq.StringArray{"interest1"},
				)

				mock.ExpectQuery("SELECT p.user_id AS user_id, p.name, p.email, COALESCE").
					WithArgs(strconv.FormatInt(userID, 10)).
					WillReturnRows(rows)
			},
			expectedParticipant: nil,
			expectedError:       errors.New("sql: Scan error on column index 0, name \"user_id\": converting NULL to int64 is unsupported"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.userID)

			participant, err := r.ReadParticipant(context.Background(), testCase.userID)

			if testCase.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedParticipant, participant)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Expectations were not met: %v", err)
			}
		})
	}
}

func TestUpdateParticipant(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	r := New(postgres.DB{DB: db})

	type mockBehavior func(participant *models.Participant)

	testTable := []struct {
		name          string
		participant   *models.Participant
		mockBehavior  mockBehavior
		expectedError error
	}{
		{
			name: "OK test",
			participant: &models.Participant{
				UserID:    1,
				Name:      "Updated Name",
				Email:     "updated@example.com",
				Interests: []string{"interest1", "interest2"},
			},
			mockBehavior: func(participant *models.Participant) {
				mock.ExpectBegin()

				mock.ExpectExec("UPDATE public.users").
					WithArgs(participant.Name, participant.Email, participant.UserID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("DELETE FROM public.interests").
					WithArgs(participant.UserID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("INSERT INTO public.interests").
					WithArgs(participant.UserID, "interest1", participant.UserID, "interest2").
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "Begin transaction error",
			participant: &models.Participant{
				UserID:    1,
				Name:      "Updated Name",
				Email:     "updated@example.com",
				Interests: []string{"interest1", "interest2"},
			},
			mockBehavior: func(participant *models.Participant) {
				mock.ExpectBegin().WillReturnError(errors.New("transaction error"))
			},
			expectedError: errors.New("transaction error"),
		},
		{
			name: "Update participant error",
			participant: &models.Participant{
				UserID:    1,
				Name:      "Updated Name",
				Email:     "updated@example.com",
				Interests: []string{"interest1", "interest2"},
			},
			mockBehavior: func(participant *models.Participant) {
				mock.ExpectBegin()

				mock.ExpectExec("UPDATE public.users").
					WithArgs(participant.Name, participant.Email, participant.UserID).
					WillReturnError(errors.New("update error"))
			},
			expectedError: errors.New("update error"),
		},
		{
			name: "No rows affected error",
			participant: &models.Participant{
				UserID:    1,
				Name:      "Updated Name",
				Email:     "updated@example.com",
				Interests: []string{"interest1", "interest2"},
			},
			mockBehavior: func(participant *models.Participant) {
				mock.ExpectBegin()

				mock.ExpectExec("UPDATE public.users").
					WithArgs(participant.Name, participant.Email, participant.UserID).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedError: models.ErrWrongUserId,
		},
		{
			name: "Delete interests error",
			participant: &models.Participant{
				UserID:    1,
				Name:      "Updated Name",
				Email:     "updated@example.com",
				Interests: []string{"interest1", "interest2"},
			},
			mockBehavior: func(participant *models.Participant) {
				mock.ExpectBegin()

				mock.ExpectExec("UPDATE public.users").
					WithArgs(participant.Name, participant.Email, participant.UserID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("DELETE FROM public.interests").
					WithArgs(participant.UserID).
					WillReturnError(errors.New("delete error"))
			},
			expectedError: errors.New("delete error"),
		},
		{
			name: "Insert interests error",
			participant: &models.Participant{
				UserID:    1,
				Name:      "Updated Name",
				Email:     "updated@example.com",
				Interests: []string{"interest1", "interest2"},
			},
			mockBehavior: func(participant *models.Participant) {
				mock.ExpectBegin()

				mock.ExpectExec("UPDATE public.users").
					WithArgs(participant.Name, participant.Email, participant.UserID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("DELETE FROM public.interests").
					WithArgs(participant.UserID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("INSERT INTO public.interests").
					WithArgs(participant.UserID, "interest1", participant.UserID, "interest2").
					WillReturnError(errors.New("insert error"))
			},
			expectedError: errors.New("insert error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.participant)

			err := r.UpdateParticipant(context.Background(), testCase.participant)

			if testCase.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Expectations were not met: %v", err)
			}
		})
	}
}

func TestSetChatStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	r := New(postgres.DB{DB: db})

	type mockBehavior func(userID int64, eventID int64, isReady bool)

	testTable := []struct {
		name          string
		userID        int64
		eventID       int64
		isReady       bool
		mockBehavior  mockBehavior
		expectedError error
	}{
		{
			name:    "OK test",
			userID:  1,
			eventID: 2,
			isReady: true,
			mockBehavior: func(userID int64, eventID int64, isReady bool) {
				mock.ExpectExec("UPDATE public.registrations").
					WithArgs(isReady, strconv.FormatInt(eventID, 10), strconv.FormatInt(userID, 10)).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: nil,
		},
		{
			name:    "Update query error",
			userID:  1,
			eventID: 2,
			isReady: true,
			mockBehavior: func(userID int64, eventID int64, isReady bool) {
				mock.ExpectExec("UPDATE public.registrations").
					WithArgs(isReady, strconv.FormatInt(eventID, 10), strconv.FormatInt(userID, 10)).
					WillReturnError(errors.New("update error"))
			},
			expectedError: errors.New("update error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.userID, testCase.eventID, testCase.isReady)

			err := r.SetChatStatus(context.Background(), testCase.userID, testCase.eventID, testCase.isReady)

			if testCase.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Expectations were not met: %v", err)
			}
		})
	}
}
