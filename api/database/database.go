package database

import (
	"context"
	"os"

	"github.com/go-redis/redis"
)

// Ctx is the context used for Redis operations.
var Ctx = context.Background()

// CreateClient creates and returns a new Redis client with the specified database number.
func CreateClient(dbNo int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DB_ADDR"), // Redis server address
		Password: os.Getenv("DB_PASS"), // Redis server password
		DB:       dbNo,                 // Database number
	})
	return rdb
}
