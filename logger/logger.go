package logger

import (
	"log"
	"log/slog"
	"os"
	"time"
)

func LogInit() {
	now := time.Now()
	f, err := os.OpenFile("log/system-"+now.Format(time.DateOnly)+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Print(err)
	}
	logger := slog.New(slog.NewJSONHandler(f, nil))
	slog.SetDefault(logger)
}
