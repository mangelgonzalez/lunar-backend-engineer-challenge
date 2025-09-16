package logger

import (
	"context"

	"go.uber.org/zap"
)

type ContextDecorator func(ctx context.Context) []zap.Field

type ContextDecorators []ContextDecorator

func (cds ContextDecorators) Decorate(ctx context.Context) []zap.Field {
	fields := make([]zap.Field, 0)
	for _, cd := range cds {
		fields = append(fields, cd(ctx)...)
	}
	return fields
}
