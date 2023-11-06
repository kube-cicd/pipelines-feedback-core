package store

import (
	"context"
	"github.com/redis/go-redis/v9"
	"os"
	"strconv"
	"time"
)

type Redis struct {
	client *redis.Client
}

func (r *Redis) Set(key, value string, ttl int) error {
	if ttl == 0 {
		ttl = 86400 * 365 * 10 // 10 years should be enough
	}
	err := r.client.Set(context.TODO(), key, value, time.Second*time.Duration(ttl)).Err()
	return err
}

func (r *Redis) Get(key string) (string, error) {
	return r.client.Get(context.TODO(), key).Result()
}

func NewRedis() *Redis {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost:6379"
	}
	redisDB := os.Getenv("REDIS_DB")
	if redisDB == "" {
		redisDB = "0" // use default DB
	}
	redisDBi, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	return &Redis{
		redis.NewClient(&redis.Options{
			Addr:     redisHost,
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       redisDBi,
		}),
	}
}
