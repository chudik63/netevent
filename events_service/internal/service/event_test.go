package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/chudik63/netevent/events_service/internal/models"
	"github.com/chudik63/netevent/events_service/internal/producer"
	"github.com/chudik63/netevent/events_service/internal/repository"
	"github.com/chudik63/netevent/events_service/internal/service/mock"
	"github.com/redis/go-redis/v9"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateRegistration(t *testing.T) {
	testTable := []struct {
		name                string
		inputUserID         int64
		inputEventID        int64
		kafkaMessage        producer.Message
		kafkaTopic          string
		mockRepositorySetup func(s *mock.MockRepository, userID int64, eventID int64)
		mockProducerSetup   func(s *mock.MockProducer, message producer.Message, topic string)
		expected            error
	}{
		{
			name:         "no error test",
			inputUserID:  1,
			inputEventID: 1,
			kafkaMessage: producer.Message{},
			kafkaTopic:   "registration",
			mockRepositorySetup: func(s *mock.MockRepository, userID int64, eventID int64) {
				s.EXPECT().ReadParticipant(gomock.Any(), userID).Return(&models.Participant{}, nil).Times(2)
				s.EXPECT().ReadEvent(gomock.Any(), eventID).Return(&models.Event{}, nil).Times(2)
				s.EXPECT().CreateRegistration(gomock.Any(), userID, eventID).Return(nil).Times(1)
			},
			mockProducerSetup: func(s *mock.MockProducer, message producer.Message, topic string) {
				s.EXPECT().Produce(gomock.Any(), message, topic).Times(1)
			},
			expected: nil,
		},
		{
			name:         "wrong event id test",
			inputUserID:  1,
			inputEventID: 99,
			kafkaMessage: producer.Message{},
			kafkaTopic:   "registration",
			mockRepositorySetup: func(s *mock.MockRepository, userID int64, eventID int64) {
				s.EXPECT().ReadParticipant(gomock.Any(), userID).Return(&models.Participant{}, nil).Times(1)
				s.EXPECT().ReadEvent(gomock.Any(), eventID).Return(&models.Event{}, models.ErrWrongEventId).Times(1)
			},
			mockProducerSetup: func(s *mock.MockProducer, message producer.Message, topic string) {
			},
			expected: models.ErrWrongEventId,
		},
		{
			name:         "wrong user id test",
			inputUserID:  99,
			inputEventID: 1,
			kafkaMessage: producer.Message{},
			kafkaTopic:   "registration",
			mockRepositorySetup: func(s *mock.MockRepository, userID int64, eventID int64) {
				s.EXPECT().ReadParticipant(gomock.Any(), userID).Return(&models.Participant{}, models.ErrWrongUserId).Times(1)
			},
			mockProducerSetup: func(s *mock.MockProducer, message producer.Message, topic string) {
			},
			expected: models.ErrWrongUserId,
		},
		{
			name:         "internal error test",
			inputUserID:  1,
			inputEventID: 1,
			kafkaMessage: producer.Message{},
			kafkaTopic:   "registration",
			mockRepositorySetup: func(s *mock.MockRepository, userID int64, eventID int64) {
				s.EXPECT().ReadParticipant(gomock.Any(), userID).Return(&models.Participant{}, nil).Times(1)
				s.EXPECT().ReadEvent(gomock.Any(), eventID).Return(&models.Event{}, nil).Times(1)
				s.EXPECT().CreateRegistration(gomock.Any(), userID, eventID).Return(errors.New("internal error")).Times(1)
			},
			mockProducerSetup: func(s *mock.MockProducer, message producer.Message, topic string) {
			},
			expected: errors.New("internal error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repository := mock.NewMockRepository(c)
			producer := mock.NewMockProducer(c)
			cache := mock.NewMockCache(c)

			testCase.mockRepositorySetup(repository, testCase.inputUserID, testCase.inputEventID)
			testCase.mockProducerSetup(producer, testCase.kafkaMessage, testCase.kafkaTopic)

			service := New(repository, cache, producer)

			err := service.CreateRegistration(context.Background(), testCase.inputUserID, testCase.inputEventID)

			assert.Equal(t, testCase.expected, err)
		})
	}
}

func TestSetChatStatus(t *testing.T) {
	testTable := []struct {
		name                string
		inputUserID         int64
		inputEventID        int64
		isReady             bool
		mockRepositorySetup func(s *mock.MockRepository, userID int64, eventID int64)
		expected            error
	}{
		{
			name:         "no error test",
			inputUserID:  1,
			inputEventID: 1,
			isReady:      true,
			mockRepositorySetup: func(s *mock.MockRepository, userID int64, eventID int64) {
				s.EXPECT().ReadParticipant(gomock.Any(), userID).Return(&models.Participant{}, nil).Times(1)
				s.EXPECT().ReadEvent(gomock.Any(), eventID).Return(&models.Event{}, nil).Times(1)
				s.EXPECT().CreateRegistration(gomock.Any(), userID, eventID).Return(models.ErrAlreadyRegistered).Times(1)
				s.EXPECT().SetChatStatus(gomock.Any(), userID, eventID, true).Return(nil).Times(1)
			},
			expected: nil,
		},
		{
			name:         "wrong user id test",
			inputUserID:  99,
			inputEventID: 1,
			isReady:      true,
			mockRepositorySetup: func(s *mock.MockRepository, userID int64, eventID int64) {
				s.EXPECT().ReadParticipant(gomock.Any(), userID).Return(nil, models.ErrWrongUserId).Times(1)
			},
			expected: models.ErrWrongUserId,
		},
		{
			name:         "wrong event id test",
			inputUserID:  1,
			inputEventID: 99,
			isReady:      true,
			mockRepositorySetup: func(s *mock.MockRepository, userID int64, eventID int64) {
				s.EXPECT().ReadParticipant(gomock.Any(), userID).Return(&models.Participant{}, nil).Times(1)
				s.EXPECT().ReadEvent(gomock.Any(), eventID).Return(nil, models.ErrWrongEventId).Times(1)
			},
			expected: models.ErrWrongEventId,
		},
		{
			name:         "user is not registered",
			inputUserID:  1,
			inputEventID: 1,
			isReady:      true,
			mockRepositorySetup: func(s *mock.MockRepository, userID int64, eventID int64) {
				s.EXPECT().ReadParticipant(gomock.Any(), userID).Return(&models.Participant{}, nil).Times(1)
				s.EXPECT().ReadEvent(gomock.Any(), eventID).Return(&models.Event{}, nil).Times(1)
				s.EXPECT().CreateRegistration(gomock.Any(), userID, eventID).Return(nil).Times(1)
			},
			expected: models.ErrRegistrationNotFound,
		},
		{
			name:         "internal error test",
			inputUserID:  1,
			inputEventID: 1,
			isReady:      true,
			mockRepositorySetup: func(s *mock.MockRepository, userID int64, eventID int64) {
				s.EXPECT().ReadParticipant(gomock.Any(), userID).Return(&models.Participant{}, nil).Times(1)
				s.EXPECT().ReadEvent(gomock.Any(), eventID).Return(&models.Event{}, nil).Times(1)
				s.EXPECT().CreateRegistration(gomock.Any(), userID, eventID).Return(models.ErrAlreadyRegistered).Times(1)
				s.EXPECT().SetChatStatus(gomock.Any(), userID, eventID, true).Return(errors.New("internal error")).Times(1)
			},
			expected: errors.New("internal error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repository := mock.NewMockRepository(c)
			producer := mock.NewMockProducer(c)
			cache := mock.NewMockCache(c)

			testCase.mockRepositorySetup(repository, testCase.inputUserID, testCase.inputEventID)

			service := New(repository, cache, producer)

			err := service.SetChatStatus(context.Background(), testCase.inputUserID, testCase.inputEventID, testCase.isReady)

			assert.Equal(t, testCase.expected, err)
		})
	}
}

func TestListRegistratedEvents(t *testing.T) {
	testTable := []struct {
		name                string
		inputUserID         int64
		mockRepositorySetup func(s *mock.MockRepository, userID int64)
		expected            []*models.Event
		expectedErr         error
	}{
		{
			name:        "no error test",
			inputUserID: 1,
			mockRepositorySetup: func(s *mock.MockRepository, userID int64) {
				s.EXPECT().ReadParticipant(gomock.Any(), userID).Return(&models.Participant{}, nil).Times(1)
				s.EXPECT().ListRegistratedEvents(gomock.Any(), userID).Return([]*models.Event{{EventID: 1}}, nil).Times(1)
			},
			expected:    []*models.Event{{EventID: 1}},
			expectedErr: nil,
		},
		{
			name:        "wrong user id test",
			inputUserID: 99,
			mockRepositorySetup: func(s *mock.MockRepository, userID int64) {
				s.EXPECT().ReadParticipant(gomock.Any(), userID).Return(nil, models.ErrWrongUserId).Times(1)
			},
			expected:    nil,
			expectedErr: models.ErrWrongUserId,
		},
		{
			name:        "internal error test",
			inputUserID: 1,
			mockRepositorySetup: func(s *mock.MockRepository, userID int64) {
				s.EXPECT().ReadParticipant(gomock.Any(), userID).Return(&models.Participant{}, nil).Times(1)
				s.EXPECT().ListRegistratedEvents(gomock.Any(), userID).Return(nil, errors.New("internal error")).Times(1)
			},
			expected:    nil,
			expectedErr: errors.New("internal error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repository := mock.NewMockRepository(c)

			testCase.mockRepositorySetup(repository, testCase.inputUserID)

			service := New(repository, nil, nil)

			result, err := service.ListRegistratedEvents(context.Background(), testCase.inputUserID)

			assert.Equal(t, testCase.expectedErr, err)
			assert.Equal(t, testCase.expected, result)
		})
	}
}

func TestListEventsByInterests(t *testing.T) {
	testTable := []struct {
		name                string
		inputUserID         int64
		inputCreds          repository.Creds
		mockRepositorySetup func(s *mock.MockRepository, userID int64, creds repository.Creds)
		mockCacheSetup      func(s *mock.MockCache, userID int64)
		expected            []*models.Event
		expectedErr         error
	}{
		{
			name:        "not cached",
			inputUserID: 1,
			inputCreds:  repository.Creds{},
			mockRepositorySetup: func(s *mock.MockRepository, userID int64, creds repository.Creds) {
				s.EXPECT().ReadParticipant(gomock.Any(), userID).Return(&models.Participant{}, nil).Times(1)
				s.EXPECT().ListEventsByInterests(gomock.Any(), userID, creds).Return([]*models.Event{{EventID: 1}}, nil).Times(1)
			},
			mockCacheSetup: func(s *mock.MockCache, userID int64) {
				s.EXPECT().Get(gomock.Any(), "1_eventsByInterests").Return(&redis.StringCmd{}).Times(1)
				s.EXPECT().Set(gomock.Any(), "1_eventsByInterests", gomock.Any(), gomock.Any()).Return(&redis.StatusCmd{}).Times(1)
			},
			expected:    []*models.Event{{EventID: 1}},
			expectedErr: nil,
		},
		{
			name:        "wrong user id test",
			inputUserID: 99,
			inputCreds:  repository.Creds{},
			mockRepositorySetup: func(s *mock.MockRepository, userID int64, creds repository.Creds) {
				s.EXPECT().ReadParticipant(gomock.Any(), userID).Return(nil, models.ErrWrongUserId).Times(1)
			},
			mockCacheSetup: func(s *mock.MockCache, userID int64) {},
			expected:       nil,
			expectedErr:    models.ErrWrongUserId,
		},
		{
			name:        "get from cache",
			inputUserID: 1,
			inputCreds:  repository.Creds{},
			mockRepositorySetup: func(s *mock.MockRepository, userID int64, creds repository.Creds) {
				s.EXPECT().ReadParticipant(gomock.Any(), userID).Return(&models.Participant{}, nil).Times(1)
			},
			mockCacheSetup: func(s *mock.MockCache, userID int64) {
				res := &redis.StringCmd{}
				data, _ := json.Marshal([]*models.Event{{EventID: 1}})
				res.SetVal(string(data))
				s.EXPECT().Get(gomock.Any(), "1_eventsByInterests").Return(res).Times(1)
			},
			expected:    []*models.Event{{EventID: 1}},
			expectedErr: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repository := mock.NewMockRepository(c)
			cache := mock.NewMockCache(c)

			testCase.mockRepositorySetup(repository, testCase.inputUserID, testCase.inputCreds)
			testCase.mockCacheSetup(cache, testCase.inputUserID)

			service := New(repository, cache, nil)

			result, err := service.ListEventsByInterests(context.Background(), testCase.inputUserID, testCase.inputCreds)

			assert.Equal(t, testCase.expectedErr, err)
			assert.Equal(t, testCase.expected, result)
		})
	}
}
