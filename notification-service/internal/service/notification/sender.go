package notification

import (
	"context"
	"time"

	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/internal/application/config"
	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/internal/domain"
	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/pkg/logger"
)

//go:generate mockery --name Mail --structname MockMail --filename mock_mail_test.go --outpkg notification_test --output .
type Mail interface {
	Send(subject string, msg domain.Notification) error
}

type Sender struct {
	repo     NotificationRepository
	mail     Mail
	interval time.Duration
	done     chan struct{}
	stopped  chan struct{}
}

func NewSender(cfg config.Sender, repo NotificationRepository, mail Mail) *Sender {
	return &Sender{
		repo:     repo,
		mail:     mail,
		interval: time.Minute * time.Duration(cfg.MinuteInterval),
		done:     make(chan struct{}),
		stopped:  make(chan struct{}),
	}
}

func (s *Sender) Run(ctx context.Context) error {
	for {
		select {
		case <-s.done:
			return nil

		default:
			notifications, err := s.repo.GetNotifications(ctx)
			if err != nil {
				logger.Default().Errorf(ctx, "failed to get notifications: %s", err)
			}

			for _, notify := range notifications {
				s.mail.Send(notify.EventName, notify)

				_, err := s.repo.DeleteNotification(ctx, notify.ID)
				if err != nil {
					logger.Default().Errorf(ctx, "failed to delete notification with id = %d: %s", notify.ID, err)
				}
			}

			time.Sleep(s.interval)
		}
	}
}

func (s *Sender) Stop(ctx context.Context) error {
	close(s.done)

	select {
	case <-s.stopped:
		break

	case <-ctx.Done():
		return NewErrStopSender()
	}

	return nil
}
