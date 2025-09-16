package test

import (
	"context"
	"lunar-backend-engineer-challenge/cmd/di"
	"sync"
)

type RedisArranger struct {
	common *di.CommonServices
}

func NewRedisArranger(common *di.CommonServices) *RedisArranger {
	return &RedisArranger{common: common}
}

func (r *RedisArranger) Arrange(ctx context.Context, wg *sync.WaitGroup) {
	status := r.common.RedisClient.FlushAllAsync(ctx)
	if err := status.Err(); err != nil {
		panic(err)
	}

	wg.Done()
}
