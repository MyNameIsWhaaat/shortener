package logger

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

func Init() {
	opts := slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	handler := slog.NewJSONHandler(os.Stdout, &opts)
	Log = slog.New(handler)
	slog.SetDefault(Log)
}

func Info(msg string, args ...any) {
	Log.Info(msg, args...)
}

func Error(msg string, args ...any) {
	Log.Error(msg, args...)
}

func Debug(msg string, args ...any) {
	Log.Debug(msg, args...)
}

func Warn(msg string, args ...any) {
	Log.Warn(msg, args...)
}
