package main

import (
	"github.com/chudik63/netevent/events_service/internal/app"
)

const (
	serviceName = "events_service"
)

func main() {
	app.Run(serviceName)
}
