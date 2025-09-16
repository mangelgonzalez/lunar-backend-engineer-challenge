package command

import (
	"context"
	"lunar-backend-engineer-challenge/pkg/bus"
)

type CommandHandler interface {
	Handle(ctx context.Context, command bus.Dto) error
}
