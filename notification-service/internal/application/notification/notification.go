package notification

import (
	"fmt"

	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/internal/config"
	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/internal/domain"
	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/pkg/mail/gmail"
)

func Start() error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	gm := gmail.New(cfg.GmailConfig)
	err = gm.Send("test", domain.Message{
		UserName:   "Denis",
		EventName:  "Test",
		EventTime:  "22.04.2025",
		EventPlace: "University",
	}, "3d.lebedzeu@gmail.com")
	if err != nil {
		return fmt.Errorf("failed to send mail: %w", err)
	}

	return nil
}
