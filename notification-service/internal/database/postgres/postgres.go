package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/internal/database"
	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/internal/domain"
)

type PostgresConfig struct {
	PostgresUserName string `env:"POSTGRES_USER"     env-default:"root"`
	PostgresPassword string `env:"POSTGRES_PASSWORD" env-default:"123"`
	PostgresDBName   string `env:"POSTGRES_DB"       env-default:"netevent"`
	PostgresHost     string `env:"POSTGRES_HOST"     env-default:"localhost"`
	PostgresPort     string `env:"POSTGRES_PORT"     env-default:"5432"`
}

type DB struct {
	db      *sqlx.DB
	queries *database.Queries
}

func NewPostgres(config PostgresConfig) (*DB, error) {
	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		config.PostgresUserName,
		config.PostgresPassword,
		config.PostgresDBName,
		config.PostgresHost,
		config.PostgresPort,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if _, err := db.Conn(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &DB{
		queries: database.New(db),
	}, nil
}

func (db *DB) GetNotifications(ctx context.Context) ([]domain.Notification, error) {
	res, err := db.queries.GetNotifications(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}

	notifies := make([]domain.Notification, 0, len(res))
	for _, v := range res {
		notifies = append(notifies, dbNotificationToGlobal(v))
	}

	return notifies, nil
}

func (db *DB) AddNotification(ctx context.Context, notify domain.Notification) (domain.Notification, error) {
	tm, err := time.Parse(time.DateTime, notify.EventTime)
	if err != nil {
		return domain.Notification{}, fmt.Errorf("failed to parse time: %w", err)
	}

	args := database.AddNotificationParams{
		UserName:   notify.UserName,
		EventName:  notify.EventName,
		EventPlace: notify.EventPlace,
		EventTime:  tm,
	}

	res, err := db.queries.AddNotification(ctx, args)
	if err != nil {
		return domain.Notification{}, fmt.Errorf("failed to add notification: %w", err)
	}

	return dbNotificationToGlobal(res), nil
}

func (db *DB) DeleteNotification(ctx context.Context, id int64) (domain.Notification, error) {
	res, err := db.queries.DeleteNotification(ctx, id)
	if err != nil {
		return domain.Notification{}, fmt.Errorf("failed to delete notification: %w", err)
	}

	return dbNotificationToGlobal(res), nil
}

func (db *DB) Close() error {
	if err := db.db.Close(); err != nil {
		return fmt.Errorf("failed to close postgres: %w", err)
	}

	return nil
}

func dbNotificationToGlobal(notify database.Notification) domain.Notification {
	return domain.Notification{
		ID:         notify.ID,
		UserName:   notify.UserName,
		EventName:  notify.EventName,
		EventPlace: notify.EventPlace,
		EventTime:  notify.EventTime.Format(time.DateTime),
	}
}
