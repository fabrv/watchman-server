package database

import (
	"sync"

	"github.com/fabrv/watchman-server/utils"
	"github.com/gofiber/storage/redis"
)

var redisLock = &sync.Mutex{}
var redisInstance *redis.Storage

func RedisStorageInstance() *redis.Storage {
	if redisInstance == nil {
		redisLock.Lock()
		defer redisLock.Unlock()

		if redisInstance == nil {
			redisInstance = redis.New(redis.Config{
				URL:   utils.GetEnv("REDIS_URL", "redis://:@127.0.0.1:6379/0"),
				Reset: false,
			})
		}
	}
	return redisInstance
}
