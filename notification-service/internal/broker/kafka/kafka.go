package kafka

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/internal/application/config"
	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/pkg/logger"
)

type Closer = func() error

type Kafka struct {
	group   sarama.ConsumerGroup
	handler sarama.ConsumerGroupHandler
	topics  []string
	done    chan struct{}
	stopped chan struct{}
}

func New(cfg config.Kafka, handler sarama.ConsumerGroupHandler) (*Kafka, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

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

	for {
		select {
		case <-k.done:
			close(k.stopped)

			return nil

		default:
			err := k.group.Consume(ctx, k.topics, k.handler)
			if err != nil {
				logger.Default().Errorf(ctx, "failed to consume messages: %s", err)
			}
		}
	}
}

func (k *Kafka) Stop(ctx context.Context) error {
	close(k.done)

	select {
	case <-k.stopped:
		break

	case <-ctx.Done():
		return NewErrStopKafkaConsumer()
	}

	if err := k.group.Close(); err != nil {
		return fmt.Errorf("failed to close consumer group: %w", err)
	}

	return nil
}
