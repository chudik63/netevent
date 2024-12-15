-- name: GetNotifications :many
SELECT 
	*
FROM notifications
WHERE AGE(NOW(), event_time) <= INTERVAL '1 day';

-- name: AddNotification :one
INSERT INTO notifications(user_name, user_email, event_name, event_place, event_time)
VALUES (sqlc.arg(user_name), sqlc.arg(user_email), sqlc.arg(event_name), sqlc.arg(event_place), sqlc.arg(event_time))
RETURNING *;

-- name: DeleteNotification :one
DELETE FROM notifications
WHERE id = sqlc.arg(id)
RETURNING *;