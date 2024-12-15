package notification

import (
	"context"
	"errors"
	"fmt"

	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/internal/application/config"
	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/internal/broker/kafka"
	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/internal/database"
	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/internal/service/notification"
	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/pkg/mail/gmail"
	"golang.org/x/sync/errgroup"
)

type Closer = func() error

type Application struct {
	kafka   *kafka.Kafka
	sender  *notification.Sender
	closers []Closer
}

func New() *Application {
	return &Application{
		closers: make([]Closer, 0),
	}
}

func (a *Application) Initialize(ctx context.Context, cfg *config.Config) error {
	db, err := database.NewAdapter(ctx, cfg.Database.SQL)
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	a.closers = append(a.closers, db.Close)

	parser := notification.NewParser(db)

	kfk, err := kafka.New(cfg.Kafka, parser)
	if err != nil {
		return fmt.Errorf("kafka.New(): %w", err)
	}

	a.kafka = kfk

	mail := gmail.New(cfg.Mail)
	a.sender = notification.NewSender(cfg.Sender, db, mail)

	return nil
}

func (a *Application) Run(ctx context.Context) error {
	eg := errgroup.Group{}

	eg.Go(func() error {
		if err := a.sender.Run(ctx); err != nil {
			return fmt.Errorf("error running sender: %w", err)
		}

		return nil
	})

	eg.Go(func() error {
		if err := a.kafka.Run(ctx); err != nil {
			return fmt.Errorf("error running kafka klient: %w", err)
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("eg.Wait(): %w", err)
	}

	return nil
}

func (a *Application) Stop(ctx context.Context) error {
	errs := make([]error, 0)

	if err := a.kafka.Stop(ctx); err != nil {
		errs = append(errs, fmt.Errorf("a.kafka.Stop(ctx): %w", err))
	}

	if err := a.sender.Stop(ctx); err != nil {
		errs = append(errs, fmt.Errorf("a.sender.Stop(ctx): %w", err))
	}

	for _, closer := range a.closers {
		if err := closer(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return fmt.Errorf("failed to stop: %w", errors.Join(errs...))
	}

	return nil
}
