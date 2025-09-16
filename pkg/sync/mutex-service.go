package sync

import (
	"context"
	"lunar-backend-engineer-challenge/pkg/logger"
	"lunar-backend-engineer-challenge/pkg/utils"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const mutexName = "mutex-identity-sync"

type MutexService interface {
	Mutex(ctx context.Context, key string, fn func() (interface{}, error)) (interface{}, error)
}

type RedisMutexService struct {
	sync   *redsync.Redsync
	logger logger.Logger
}

func NewRedisMutexService(redisClient *redis.Client, logger logger.Logger) *RedisMutexService {
	pool := goredis.NewPool(redisClient)

	return &RedisMutexService{sync: redsync.New(pool), logger: logger}
}

func (rm *RedisMutexService) Mutex(ctx context.Context, key string, fn func() (interface{}, error)) (interface{}, error) {
	mutex := rm.sync.NewMutex(
		mutexName+key,
		redsync.WithExpiry(29*time.Second),
		redsync.WithRetryDelay(25*time.Millisecond),
		redsync.WithTimeoutFactor(0.05),
	)

	if _, err := utils.RetryFunc(func() (interface{}, error) {
		return nil, rm.adquireLock(ctx, mutex)
	}, 4); err != nil {
		rm.logger.Error("error locking mutex sync", zap.Error(err), zap.String("mutex_key", mutex.Name()))
		return nil, NewErrorLockMutexKey(key, err)
	}

	result, err := fn()
	if _, err := utils.RetryFunc(func() (interface{}, error) {
		return nil, rm.releaseLock(ctx, mutex)
	}, 4); err != nil {
		rm.logger.Error("error unlocking mutex sync", zap.Error(err), zap.String("mutex_key", mutex.Name()))
		return nil, NewErrorReleaseLockMutexKey(key, err)
	}

	return result, err
}

func (rm *RedisMutexService) releaseLock(ctx context.Context, mutex *redsync.Mutex) error {
	if ok, err := mutex.UnlockContext(ctx); !ok || err != nil {
		rm.logger.Warn("error unlocking mutex sync. Retrying", zap.Error(err), zap.String("mutex_key", mutex.Name()))
		if err != nil {
			return err
		}

		return errors.New("redis mutex invalid status when unlocking")
	}

	return nil
}

func (rm *RedisMutexService) adquireLock(ctx context.Context, mutex *redsync.Mutex) error {

	if err := mutex.LockContext(ctx); err != nil {
		rm.logger.Warn("error locking mutex sync. Retrying", zap.Error(err), zap.String("mutex_key", mutex.Name()))
		return err
	}

	return nil
}
