package notification

import (
	"context"
	"time"

	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/internal/application/config"
	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/internal/domain"
	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/pkg/logger"
)

//go:generate mockery --name NotificationRepository  --structname MockNotificationRepository --filename mock_notification_repository_test.go --outpkg notification_test --output .
type NotificationRepository interface {
	GetNearestNotifications(ctx context.Context) ([]domain.Notification, error)
	AddNotification(ctx context.Context, notify domain.Notification) (domain.Notification, error)
	DeleteNotification(ctx context.Context, id int64) (domain.Notification, error)
}

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
		interval: time.Second * time.Duration(cfg.SecondInterval),
		done:     make(chan struct{}),
		stopped:  make(chan struct{}),
	}
}

func (s *Sender) Run(ctx context.Context) error {
	defer close(s.stopped)

	for {
		select {
		case <-s.done:
			return nil

		case <-time.After(s.interval):
			notifications, err := s.repo.GetNearestNotifications(ctx)
			if err != nil {
				logger.GetLoggerFromCtx(ctx).Errorf(ctx, "failed to get notifications: %s", err)
				break
			}

			for _, notify := range notifications {
				if err := s.mail.Send(notify.EventName, notify); err != nil {
					logger.GetLoggerFromCtx(ctx).Errorf(ctx, "failed to send notification: %s", err)
					break
				}

				logger.GetLoggerFromCtx(ctx).Infof(ctx, "send notification to %q", notify.UserEmail)

				_, err := s.repo.DeleteNotification(ctx, notify.ID)
				if err != nil {
					logger.GetLoggerFromCtx(ctx).Errorf(ctx, "failed to delete notification with id = %d: %s", notify.ID, err)
				}
			}
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
