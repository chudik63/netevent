package gmail

import (
	"bytes"
	"fmt"
	"html/template"

	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/internal/domain"
	"gopkg.in/gomail.v2"
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
	dialer       *gomail.Dialer
	from         string
}

func New(cfg GmailConfig) *Gmail {
	return &Gmail{
		mailTemplate: template.Must(template.ParseFiles(templatePath)),
		dialer:       gomail.NewDialer(cfg.GmailHost, cfg.GmailPort, cfg.GmailUsername, cfg.GmailPassword),
		from:         cfg.GmailUsername,
	}
}

func (g *Gmail) Send(subject string, msg domain.Message, to string) error {
	var body bytes.Buffer

	if err := g.mailTemplate.Execute(&body, msg); err != nil {
		return fmt.Errorf("failed to execute mail template: %w", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", g.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body.String())

	if err := g.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send mail: %w", err)
	}

	return nil
}
