package store

import (
	"context"
	"github.com/pkg/errors"
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
	fetch := r.client.Get(context.TODO(), key)
	if fetch.Err() == redis.Nil {
		return "", errors.New(ErrNotFound)
	}
	return fetch.Result()
}

func (r *Redis) CanHandle(adapterName string) bool {
	return adapterName == r.GetImplementationName()
}
func (r *Redis) GetImplementationName() string {
	return "redis"
}

func (r *Redis) Initialize() error {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost:6379"
	}
	redisDB := os.Getenv("REDIS_DB")
	if redisDB == "" {
		redisDB = "0" // use default DB
	}
	redisDBi, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	r.client = redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redisDBi,
	})
	return nil
}

func NewRedis() *Redis {
	return &Redis{}
}
