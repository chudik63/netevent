package kafka_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/IBM/sarama"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/internal/application/config"
	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/internal/broker/kafka"
	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/internal/database"
	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/internal/domain"
	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/internal/service/notification"
)

func TestKafka(t *testing.T) {
	ctx := context.Background()
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

	parser := notification.NewParser(db)

	kfk, err := kafka.New(cfg.Kafka, parser)
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

	producer, err := sarama.NewSyncProducer([]string{fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)}, producerConfig)
	require.NoError(t, err, "failed to create kafka producer")

	defer func() {
		err := producer.Close()
		assert.NoError(t, err, "failed to close kafka producer")
	}()

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

	testMessage := "invalid massage"
	msg := &sarama.ProducerMessage{
		Topic: cfg.Kafka.Topic,
		Value: sarama.StringEncoder(testMessage),
	}

	_, _, err = producer.SendMessage(msg)
	require.NoError(t, err, "failed to send message")

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
