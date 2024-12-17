package domain

type Notification struct {
	ID         int64
	UserName   string `json:"user_name"`
	UserEmail  string `json:"user_email"`
	EventName  string `json:"event_name"`
	EventTime  string `json:"event_time"`
	EventPlace string `json:"event_place"`
}
