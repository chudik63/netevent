package gmail

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/internal/domain"
)

const templatePath = "./templates/mail.html"

type GmailConfig struct {
	GmailHost     string `env:"GMAIL_HOST"`
	GmailPort     int    `env:"GMAIL_PORT" env-default:"587"`
	GmailUsername string `env:"GMAIL_USERNAME"`
	GmailPassword string `env:"GMAIL_PASSWORD"`
}

type Gmail struct {
	mailTemplate *template.Template
	auth         smtp.Auth
	address      string
	from         string
}

func New(cfg GmailConfig) *Gmail {
	return &Gmail{
		mailTemplate: template.Must(template.ParseFiles(templatePath)),
		auth:         smtp.PlainAuth("", cfg.GmailUsername, cfg.GmailPassword, cfg.GmailHost),
		address:      fmt.Sprintf("%s:%d", cfg.GmailHost, cfg.GmailPort),
		from:         cfg.GmailUsername,
	}
}

func (g *Gmail) Send(subject string, msg domain.Message, to string) error {
	var body bytes.Buffer

	if err := g.mailTemplate.Execute(&body, msg); err != nil {
		return fmt.Errorf("failed to execute mail template: %w", err)
	}

	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"
	mail := "Subject: " + subject + "\n" +
		headers + "\n\n" +
		body.String()

	if err := smtp.SendMail(g.address, g.auth, g.from, []string{to}, []byte(mail)); err != nil {
		return fmt.Errorf("failed to send mail: %w", err)
	}

	return nil
}
