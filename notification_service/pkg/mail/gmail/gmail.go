package gmail

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/chudik63/netevent/notification_service/internal/application/config"
	"github.com/chudik63/netevent/notification_service/internal/domain"
)

const templatePath = "./templates/mail.html"

type Gmail struct {
	mailTemplate *template.Template
	auth         smtp.Auth
	address      string
	from         string
}

func New(cfg config.Mail) *Gmail {
	return &Gmail{
		mailTemplate: template.Must(template.ParseFiles(templatePath)),
		auth:         smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host),
		address:      fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		from:         cfg.Username,
	}
}

func (g *Gmail) Send(subject string, msg domain.Notification) error {
	var body bytes.Buffer

	if err := g.mailTemplate.Execute(&body, msg); err != nil {
		return fmt.Errorf("failed to execute mail template: %w", err)
	}

	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"
	mail := "Subject: " + subject + "\n" +
		headers + "\n\n" +
		body.String()

	if err := smtp.SendMail(g.address, g.auth, g.from, []string{msg.UserEmail}, []byte(mail)); err != nil {
		return fmt.Errorf("failed to send mail: %w", err)
	}

	return nil
}
