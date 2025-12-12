package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	l *slog.Logger
}

func NewLogger() *Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	return &Logger{l: slog.New(handler)}
}

func (lg *Logger) Info(msg string, args ...any) {
	lg.l.Info(msg, args...)
}

func (lg *Logger) Error(msg string, args ...any) {
	lg.l.Error(msg, args...)
}

func (lg *Logger) Debug(msg string, args ...any) {
	lg.l.Debug(msg, args...)
}


