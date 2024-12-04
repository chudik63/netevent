package main

import (
	"fmt"
	"log/slog"
	"os"

	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/internal/application/notification"
)

func main() {
	if err := notification.Start(); err != nil {
		slog.Error(fmt.Sprintf("notification.Start(): %s", err))
		os.Exit(1)
	}
}
