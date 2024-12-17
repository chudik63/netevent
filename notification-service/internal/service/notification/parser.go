package notification

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/internal/domain"
	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/pkg/logger"
)

//go:generate mockery --name NotificationRepository  --structname MockNotificationRepository --filename mock_notification_repository_test.go --outpkg notification_test --output .
type NotificationRepository interface {
	GetNearestNotifications(ctx context.Context) ([]domain.Notification, error)
	AddNotification(ctx context.Context, notify domain.Notification) (domain.Notification, error)
	DeleteNotification(ctx context.Context, id int64) (domain.Notification, error)
}

type Parser struct {
	repo NotificationRepository
}

func NewParser(repo NotificationRepository) *Parser {
	return &Parser{
		repo: repo,
	}
}

// To implement sarama.ConsumerGroupHandler interface
func (s *Parser) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// To implement sarama.ConsumerGroupHandler interface
func (s *Parser) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (s *Parser) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ctx := sess.Context()

	for msg := range claim.Messages() {
		var notification domain.Notification

		if err := json.Unmarshal(msg.Value, &notification); err != nil {
			logger.GetLoggerFromCtx(ctx).Errorf(ctx, "failed to unmarshal message %q: %s", msg.Value, err)
			sess.MarkMessage(msg, "")
			continue
		}

		_, err := s.repo.AddNotification(ctx, notification)
		if err != nil {
			logger.GetLoggerFromCtx(ctx).Errorf(ctx, "failed add notification: %s", err)
			continue
		}

		sess.MarkMessage(msg, "")
		logger.GetLoggerFromCtx(ctx).Infof(ctx, "add notification to %q", notification.UserEmail)
	}

	return nil
}
