CREATE TABLE notifications
(
  id          BIGINT                   GENERATED ALWAYS AS IDENTITY UNIQUE,
  user_name   TEXT                     NOT NULL,
  user_email  TEXT                     NOT NULL,
  event_name  TEXT                     NOT NULL,
  event_place TEXT                     NOT NULL,
  event_time  TIMESTAMP WITH TIME ZONE NOT NULL,
  PRIMARY KEY (id)
);