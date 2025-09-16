package query

import (
	"context"
	"lunar-backend-engineer-challenge/pkg/bus"
)

type QueryHandler interface {
	Handle(ctx context.Context, query bus.Dto) (interface{}, error)
}
