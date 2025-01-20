package producer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/chudik63/netevent/events_service/pkg/logger"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

const (
	flushTimeout      = 500 * time.Millisecond
	RegistrationTopic = "registration"
)

type Producer struct {
	producer sarama.SyncProducer
	logger   logger.Logger
}

func New(ctx context.Context, address []string) (*Producer, error) {
	l := logger.GetLoggerFromCtx(ctx)

	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(address, config)
	if err != nil {
		l.Fatal(ctx, "failed to create Kafka producer", zap.String("err", err.Error()))
		return nil, err
	}

	return &Producer{
		producer: producer,
		logger:   l,
	}, nil
}

func (p *Producer) Produce(ctx context.Context, message Message, topic string) {
	messageJSON, err := json.Marshal(message)
	if err != nil {
		p.logger.Error(context.Background(), "kafka: failed to marshal message", zap.String("err", err.Error()))
		return
	}

	kafkaMsg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(messageJSON),
	}

	const maxRetries = 5
	for attempt := 1; attempt <= maxRetries; attempt++ {
		partition, offset, err := p.producer.SendMessage(kafkaMsg)
		if err != nil {
			p.logger.Error(ctx, "kafka: failed to send message",
				zap.String("err", err.Error()),
				zap.Int("attempt", attempt),
			)

			if attempt == maxRetries {
				p.logger.Error(ctx, "kafka: all retry attempts failed",
					zap.String("topic", topic),
					zap.String("err", err.Error()),
				)
				return
			}

			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}

		p.logger.Info(ctx, "kafka: message sent",
			zap.String("topic", topic),
			zap.Int32("partition", partition),
			zap.Int64("offset", offset),
		)
		return
	}
}
func (p *Producer) Close() {
	if err := p.producer.Close(); err != nil {
		p.logger.Error(context.Background(), "kafka: failed to close producer", zap.String("err", err.Error()))
	}
}
