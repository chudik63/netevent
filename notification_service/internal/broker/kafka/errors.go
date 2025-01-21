package kafka

type ErrStopKafkaConsumer struct{}

func NewErrStopKafkaConsumer() error {
	return ErrStopKafkaConsumer{}
}

func (e ErrStopKafkaConsumer) Error() string {
	return "failed to stop kafka consumer"
}
