package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/fachrunwira/basic-go-api-template/config"
	"github.com/fachrunwira/basic-go-api-template/lib/logger"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	Ctx         = context.Background()
	redisLogger = *logger.SetLogger("./storage/log/cache.log")
	prefix      = getEnv("APP_NAME", "go_redis")
)

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultVal
}

func InitRedis() {
	redis_cfg := config.LoadRedisConfig()

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redis_cfg.Address, redis_cfg.Port),
		Password: redis_cfg.Password,
		DB:       redis_cfg.Database,
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		redisLogger.Printf("Failed to connect to redis: %v", err)
	}
}

func Set(key string, value interface{}, ttl *int) error {
	var expired time.Duration
	if ttl != nil {
		expired = time.Duration(*ttl)
	} else {
		expired = 0
	}

	redisKey := fmt.Sprintf("%s_%s", prefix, key)

	err := RedisClient.Set(Ctx, redisKey, value, expired).Err()
	if err != nil {
		redisLogger.Printf("Redis failed to SET: %v", err)
	}

	return err
}

func Get(key string, dest interface{}) error {
	redisKey := fmt.Sprintf("%s_%s", prefix, key)

	val, err := RedisClient.Get(Ctx, redisKey).Result()
	if err != nil {
		if err == redis.Nil {
			redisLogger.Printf("Key not found: %v", err)
			return nil
		}

		redisLogger.Printf("Redis Get Error: %v", err)
		return err
	}

	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		redisLogger.Printf("JSON Unmarshal error on Key %s: %v", key, err)
		return err
	}

	return nil
}
