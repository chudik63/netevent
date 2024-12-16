package producer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/logger"
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

	partition, offset, err := p.producer.SendMessage(kafkaMsg)
	if err != nil {
		p.logger.Error(ctx, "kafka: failed to send message", zap.String("err", err.Error()))
		return
	}

	p.logger.Info(ctx, "kafka: message sent",
		zap.String("topic", topic),
		zap.Int32("partition", partition),
		zap.Int64("offset", offset),
	)
}

func (p *Producer) Close() {
	if err := p.producer.Close(); err != nil {
		p.logger.Error(context.Background(), "kafka: failed to close producer", zap.String("err", err.Error()))
	}
}
