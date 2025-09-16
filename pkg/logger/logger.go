package logger

import (
	"context"

	"go.uber.org/zap"
)

type Logger interface {
	Error(message string, fields ...zap.Field)
	Fatal(message string, fields ...zap.Field)
	Warn(message string, fields ...zap.Field)
	Info(message string, fields ...zap.Field)
}

type ContextLogger interface {
	Logger
	ErrorContext(ctx context.Context, message string, fields ...zap.Field)
	FatalContext(ctx context.Context, message string, fields ...zap.Field)
	WarnContext(ctx context.Context, message string, fields ...zap.Field)
	InfoContext(ctx context.Context, message string, fields ...zap.Field)
}
