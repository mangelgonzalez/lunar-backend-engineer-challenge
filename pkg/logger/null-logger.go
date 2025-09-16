package logger

import (
	"context"

	"go.uber.org/zap"
)

type NullLogger struct {
}

func NewNullLogger() Logger {
	return &NullLogger{}
}

func (n NullLogger) Error(_ string, _ ...zap.Field) {
}
func (n NullLogger) Fatal(_ string, _ ...zap.Field) {
}
func (n NullLogger) Warn(_ string, _ ...zap.Field) {
}
func (n NullLogger) Info(_ string, _ ...zap.Field) {
}

type NullContextLogger struct {
	Logger
}

func NewNullContextLogger() ContextLogger {
	return &NullContextLogger{
		Logger: NewNullLogger(),
	}
}

func (nl NullContextLogger) FatalContext(_ context.Context, _ string, _ ...zap.Field) {}

func (nl NullContextLogger) ErrorContext(_ context.Context, _ string, _ ...zap.Field) {}

func (nl NullContextLogger) WarnContext(_ context.Context, _ string, _ ...zap.Field) {}

func (nl NullContextLogger) InfoContext(_ context.Context, _ string, _ ...zap.Field) {}
