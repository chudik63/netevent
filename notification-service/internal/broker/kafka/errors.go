package kafka

import "fmt"

type ErrStopKafkaConsumer struct{}

func NewErrStopKafkaConsumer() error {
	return ErrStopKafkaConsumer{}
}

func (e ErrStopKafkaConsumer) Error() string {
	return fmt.Sprintf("failed to stop kafka consumer")
}
