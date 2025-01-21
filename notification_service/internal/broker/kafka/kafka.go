package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/chudik63/netevent/notification_service/internal/application/config"
	"github.com/chudik63/netevent/notification_service/internal/domain"
	"github.com/chudik63/netevent/notification_service/pkg/logger"

	"github.com/IBM/sarama"
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

type Kafka struct {
	repo     NotificationRepository
	consumer sarama.PartitionConsumer
}

func New(cfg config.Kafka, repo NotificationRepository) (*Kafka, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Retry.Backoff = 5 * time.Second

	consumer, err := sarama.NewConsumer([]string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)}, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create new kafka consumer: %w", err)
	}

	partitionConsumer, err := consumer.ConsumePartition(cfg.Topic, cfg.Partition, sarama.OffsetNewest)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka partition consumer: %w", err)
	}

	return &Kafka{
		repo:     repo,
		consumer: partitionConsumer,
	}, nil
}

func (k *Kafka) Run(ctx context.Context) error {
	for msg := range k.consumer.Messages() {
		var notification domain.Notification

		if err := json.Unmarshal(msg.Value, &notification); err != nil {
			logger.GetLoggerFromCtx(ctx).Errorf(ctx, "failed to unmarshal message %q: %s", msg.Value, err)
			continue
		}

		_, err := k.repo.AddNotification(ctx, notification)
		if err != nil {
			logger.GetLoggerFromCtx(ctx).Errorf(ctx, "failed add notification: %s", err)
			continue
		}

		logger.GetLoggerFromCtx(ctx).Infof(ctx, "add notification to %q", notification.UserEmail)
	}

	return nil
}

func (k *Kafka) Stop(ctx context.Context) error {
	if err := k.consumer.Close(); err != nil {
		return fmt.Errorf("failed to close consumer group: %w", err)
	}

	return nil
}
