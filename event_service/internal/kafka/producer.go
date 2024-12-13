package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/logger"
	"go.uber.org/zap"
)

const (
	flushTimeout      = 500
	RegistrationTopic = "registration"
)

type Producer struct {
	producer *kafka.Producer
	logger   logger.Logger
}

func NewProducer(ctx context.Context, adress string) (*Producer, error) {
	l := logger.GetLoggerFromCtx(ctx)

	conf := &kafka.ConfigMap{
		"bootstrap.servers": adress,
	}

	p, err := kafka.NewProducer(conf)
	if err != nil {
		l.Fatal(ctx, "failed to create kafka producer", zap.String("err", err.Error()))
	}

	return &Producer{
		producer: p,
		logger:   l,
	}, nil
}

func (p *Producer) Produce(ctx context.Context, message Message, topic string) {
	messageJSON, err := json.Marshal(message)
	if err != nil {
		p.logger.Error(context.Background(), "kafka: failed to read message", zap.String("err", err.Error()))
		return
	}

	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(messageJSON),
	}

	kafkaChan := make(chan kafka.Event)

	go func(ctx context.Context) {
		defer close(kafkaChan)
		if err := p.producer.Produce(kafkaMsg, kafkaChan); err != nil {
			kafkaChan <- kafka.NewError(kafka.ErrUnknown, err.Error(), false)
		}

		select {
		case e := <-kafkaChan:
			switch ev := e.(type) {
			case *kafka.Error:
				p.logger.Error(ctx, "kafka: failed to send message", zap.String("err", ev.Error()))
			default:
				p.logger.Error(ctx, "kafka: unknown event type", zap.String("event", fmt.Sprintf("%T", e)))
			}
		}
	}(ctx)
}

func (p *Producer) Close() {
	p.producer.Flush(flushTimeout)
	p.producer.Close()
}
