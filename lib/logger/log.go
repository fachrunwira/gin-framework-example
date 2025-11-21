package logger

import (
	"log/slog"

	"gopkg.in/natefinch/lumberjack.v2"
)

func New(filename string, level slog.Level) *slog.Logger {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    10,   // Max file size before rotation
		MaxBackups: 3,    // Max old log files to retain
		MaxAge:     28,   // Max days to retain a log files
		Compress:   true, // Compress old files
	}

	handler := slog.NewTextHandler(lumberjackLogger, &slog.HandlerOptions{
		Level: level,
	})

	return slog.New(handler)
	// return log.New(lumberjackLogger, "", log.LstdFlags|log.Lshortfile)
}
