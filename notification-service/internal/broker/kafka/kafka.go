package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/internal/application/config"
	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/pkg/logger"
)

type Kafka struct {
	group   sarama.ConsumerGroup
	handler sarama.ConsumerGroupHandler
	topics  []string
}

func New(cfg config.Kafka, handler sarama.ConsumerGroupHandler) (*Kafka, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.MaxWaitTime = 100 * time.Millisecond
	config.Consumer.MaxProcessingTime = 100 * time.Millisecond

	group, err := sarama.NewConsumerGroup([]string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)}, cfg.Group, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new kafka consumer: %w", err)
	}

	return &Kafka{
		group:   group,
		handler: handler,
		topics:  []string{cfg.Topic},
	}, nil
}

func (k *Kafka) Run(ctx context.Context) error {
	// Track errors
	go func() {
		for err := range k.group.Errors() {
			logger.Default().Errorf(ctx, "kafka group error: %s", err)
		}
	}()

	err := k.group.Consume(ctx, k.topics, k.handler)
	if err != nil {
		return fmt.Errorf("failed to consume messages: %w", err)
	}

	return nil
}

func (k *Kafka) Stop(ctx context.Context) error {
	if err := k.group.Close(); err != nil {
		return fmt.Errorf("failed to close consumer group: %w", err)
	}

	return nil
}
