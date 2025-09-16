package redis

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func InitRedisClient(redisHost string, redisPort int, poolSize int, idleTimeout int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%s:%d", redisHost, redisPort),
		IdleTimeout: time.Second * time.Duration(idleTimeout),
		PoolSize:    poolSize,
	})
}
