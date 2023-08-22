package repository

import "github.com/go-redis/redis" // Importing the Go Redis package

// Repository represents a structure that holds Redis clients for rate limiting and short URL storage.
type Repository struct {
	RateLimitDBClient *redis.Client // Redis client for rate limiting
	ShortUrlDBClient  *redis.Client // Redis client for short URL storage
}
