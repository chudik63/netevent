package notification

import "fmt"

type ErrStopSender struct{}

func NewErrStopSender() error {
	return ErrStopSender{}
}

func (e ErrStopSender) Error() string {
	return fmt.Sprintf("failed to stop sender")
}
