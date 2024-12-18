package producer

type Message struct {
	UserEmail  string `json:"user_email"`
	UserName   string `json:"user_name"`
	EventName  string `json:"event_name"`
	EventTime  string `json:"event_time"`
	EventPlace string `json:"event_place"`
}
