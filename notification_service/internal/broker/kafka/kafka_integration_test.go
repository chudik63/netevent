//go:build integration
// +build integration

package kafka_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/chudik63/netevent/notification-service/internal/application/config"
	"github.com/chudik63/netevent/notification-service/internal/broker/kafka"
	"github.com/chudik63/netevent/notification-service/internal/database"
	"github.com/chudik63/netevent/notification-service/internal/domain"
	"github.com/chudik63/netevent/notification-service/pkg/logger"

	"github.com/IBM/sarama"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKafka(t *testing.T) {
	lg := logger.New(os.Stdout, slog.LevelInfo, "test-notification-service")
	ctx := logger.CtxWithLogger(context.Background(), lg)
	cfg := config.Config{}

	err := cleanenv.ReadConfig("../../../test/.env", &cfg)
	require.NoError(t, err, "failed to read config")

	db, err := database.NewAdapter(ctx, cfg.Database.SQL)
	require.NoError(t, err, "failed to create db adapter")

	defer func() {
		err := db.Close()
		assert.NoError(t, err, "failed to close db")
	}()

	defer func() {
		dbNotifications, err := db.GetAllNotifications(ctx)
		require.NoError(t, err, "failed to get all notifications")

		for _, notify := range dbNotifications {
			_, err := db.DeleteNotification(ctx, notify.ID)
			assert.NoError(t, err, "failed to delete notification")
		}
	}()

	kfk, err := kafka.New(cfg.Kafka, db)
	require.NoError(t, err, "failed to create kafka")

	go func() {
		err := kfk.Run(ctx)
		assert.NoError(t, err, "failed to run kafka")
	}()

	defer func() {
		err := kfk.Stop(ctx)
		assert.NoError(t, err, "failed to stop kafka")
	}()

	producerConfig := sarama.NewConfig()
	producerConfig.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{"localhost:9091", "localhost:9092", "localhost:9093"}, producerConfig)
	require.NoError(t, err, "failed to create kafka producer")

	notifications := []domain.Notification{
		{
			UserName:   "Ivan Ivanov",
			UserEmail:  "ivan.ivanov@example.com",
			EventName:  "Software Development Conference",
			EventTime:  time.Now().Add(12 * time.Hour).Format(time.DateTime),
			EventPlace: "Moscow, Primer Street, Building 10",
		},
		{
			UserName:   "Petr Petrov",
			UserEmail:  "petr.petrov@example.com",
			EventName:  "Go Development Meetup",
			EventTime:  time.Now().Add(6 * time.Hour).Format(time.DateTime),
			EventPlace: "Saint Petersburg, Nevsky Avenue, Building 20",
		},
		{
			UserName:   "Alexey Sidorov",
			UserEmail:  "sidorov.alex@example.com",
			EventName:  "Machine Learning Hackathon",
			EventTime:  "2024-07-01 10:00:00",
			EventPlace: "Yekaterinburg, Academic Street, Building 5",
		},
	}

	for _, notify := range notifications {
		notifyJSON, err := json.Marshal(notify)
		require.NoError(t, err, "failed to marshal notification")

		msg := &sarama.ProducerMessage{
			Topic: cfg.Kafka.Topic,
			Value: sarama.StringEncoder(notifyJSON),
		}

		_, _, err = producer.SendMessage(msg)
		require.NoError(t, err, "failed to send message")
	}

	invalidMessage := "invalid massage"
	msg := &sarama.ProducerMessage{
		Topic: cfg.Kafka.Topic,
		Value: sarama.StringEncoder(invalidMessage),
	}

	_, _, err = producer.SendMessage(msg)
	require.NoError(t, err, "failed to send message")

	err = producer.Close()
	assert.NoError(t, err, "failed to close kafka producer")

	// waiting for work
	time.Sleep(5 * time.Second)

	nearestNotification, err := db.GetNearestNotifications(ctx)
	require.NoError(t, err)

	// id may differ
	for i := range nearestNotification {
		nearestNotification[i].ID = 0
	}

	assert.ElementsMatch(t, notifications[:2], nearestNotification, "nearest notifications should match")

	dbNotifications, err := db.GetAllNotifications(ctx)
	require.NoError(t, err, "failed to get all notifications")

	// id may differ
	for i := range dbNotifications {
		dbNotifications[i].ID = 0
	}

	assert.ElementsMatch(t, notifications, dbNotifications, "notifications should match")
}
