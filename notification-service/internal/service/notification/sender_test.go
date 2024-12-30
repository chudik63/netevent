package notification_test

import (
	context "context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/chudik63/netevent/notification-service/internal/application/config"
	"github.com/chudik63/netevent/notification-service/internal/domain"
	"github.com/chudik63/netevent/notification-service/internal/service/notification"
	"github.com/chudik63/netevent/notification-service/pkg/logger"

	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func TestRun(t *testing.T) {
	mockRepo := NewMockNotificationRepository(t)
	mockRepo.On("GetNearestNotifications", mock.Anything).Return([]domain.Notification{
		{
			ID:         1,
			UserName:   "name",
			UserEmail:  "main",
			EventName:  "Event 1",
			EventPlace: "place",
			EventTime:  "2006-01-02 15:04:05",
		},
		{
			ID:         2,
			UserName:   "name",
			UserEmail:  "main",
			EventName:  "Event 2",
			EventPlace: "place",
			EventTime:  "2006-01-02 15:04:05",
		},
	}, nil)
	mockRepo.On("DeleteNotification", mock.Anything, int64(1)).Return(domain.Notification{
		ID:         1,
		UserName:   "name",
		UserEmail:  "main",
		EventName:  "event",
		EventPlace: "place",
		EventTime:  "2006-01-02 15:04:05",
	}, nil)
	mockRepo.On("DeleteNotification", mock.Anything, int64(2)).Return(domain.Notification{
		ID:         2,
		UserName:   "name",
		UserEmail:  "main",
		EventName:  "event",
		EventPlace: "place",
		EventTime:  "2006-01-02 15:04:05",
	}, nil)

	mockMail := NewMockMail(t)
	mockMail.On("Send", "Event 1", mock.Anything).Return(nil)
	mockMail.On("Send", "Event 2", mock.Anything).Return(nil)

	cfg := config.Sender{SecondInterval: 1}

	sender := notification.NewSender(cfg, mockRepo, mockMail)
	lg := logger.New(os.Stdout, slog.LevelInfo, "test-notification-service")
	ctx := logger.CtxWithLogger(context.Background(), lg)

	go func() {
		err := sender.Run(ctx)
		assert.NoError(t, err)
	}()

	// waiting for work
	time.Sleep(1800 * time.Millisecond)

	err := sender.Stop(ctx)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
	mockMail.AssertExpectations(t)
}
