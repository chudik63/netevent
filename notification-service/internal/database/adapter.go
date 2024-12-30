package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chudik63/netevent/notification-service/internal/application/config"
	"github.com/chudik63/netevent/notification-service/internal/domain"
)

type Closer = func() error

type DBAdapter struct {
	queries *Queries
	closers []Closer
}

func NewAdapter(ctx context.Context, cfg config.SQL) (*DBAdapter, error) {
	closers := make([]Closer, 0)

	sqlDB, err := sql.Open(cfg.Driver, cfg.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	closers = append(closers, sqlDB.Close)

	return &DBAdapter{
		queries: New(sqlDB),
		closers: closers,
	}, nil
}

func (db *DBAdapter) GetAllNotifications(ctx context.Context) ([]domain.Notification, error) {
	res, err := db.queries.GetAllNotifications(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}

	notifies := make([]domain.Notification, 0, len(res))
	for _, v := range res {
		notifies = append(notifies, dbNotificationToGlobal(v))
	}

	return notifies, nil
}

func (db *DBAdapter) GetNearestNotifications(ctx context.Context) ([]domain.Notification, error) {
	res, err := db.queries.GetNearestNotifications(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}

	notifies := make([]domain.Notification, 0, len(res))
	for _, v := range res {
		notifies = append(notifies, dbNotificationToGlobal(v))
	}

	return notifies, nil
}

func (db *DBAdapter) AddNotification(ctx context.Context, notify domain.Notification) (domain.Notification, error) {
	tm, err := time.Parse(time.DateTime, notify.EventTime)
	if err != nil {
		return domain.Notification{}, fmt.Errorf("failed to parse time: %w", err)
	}

	args := AddNotificationParams{
		UserName:   notify.UserName,
		UserEmail:  notify.UserEmail,
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

func (db *DBAdapter) DeleteNotification(ctx context.Context, id int64) (domain.Notification, error) {
	res, err := db.queries.DeleteNotification(ctx, id)
	if err != nil {
		return domain.Notification{}, fmt.Errorf("failed to delete notification: %w", err)
	}

	return dbNotificationToGlobal(res), nil
}

func (db *DBAdapter) Close() error {
	errs := make([]error, 0)

	for _, closer := range db.closers {
		if err := closer(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return fmt.Errorf("failed to close database: %w", errors.Join(errs...))
	}

	return nil
}

func dbNotificationToGlobal(notify Notification) domain.Notification {
	return domain.Notification{
		ID:         notify.ID,
		UserName:   notify.UserName,
		UserEmail:  notify.UserEmail,
		EventName:  notify.EventName,
		EventPlace: notify.EventPlace,
		EventTime:  notify.EventTime.Format(time.DateTime),
	}
}
