-- +goose Up
-- +goose StatementBegin
CREATE TABLE notifications
(
  id          BIGINT                   GENERATED ALWAYS AS IDENTITY UNIQUE,
  user_name   TEXT                     NOT NULL,
  event_name  TEXT                     NOT NULL,
  event_place TEXT                     NOT NULL,
  event_time  TIMESTAMP WITH TIME ZONE NOT NULL,
  PRIMARY KEY (id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notifications;
-- +goose StatementEnd
