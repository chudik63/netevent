package kafka

import (
	"context"
	"encoding/json"
	"errors"

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

	return &Producer{producer: p}, nil
}

func (p *Producer) Produce(message Message, topic string) error {
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(messageJSON),
	}

	kafkaChan := make(chan kafka.Event)

	go func() {
		if err := p.producer.Produce(kafkaMsg, kafkaChan); err != nil {
			kafkaChan <- kafka.NewError(kafka.ErrUnknown, err.Error(), false)
		}
	}()

	select {
	case e := <-kafkaChan:
		switch ev := e.(type) {
		case *kafka.Message:
			return nil
		case *kafka.Error:
			return ev
		default:
			return errors.New("unknown event type")
		}
	}
}

func (p *Producer) Close() {
	p.producer.Flush(flushTimeout)
	p.producer.Close()
}
