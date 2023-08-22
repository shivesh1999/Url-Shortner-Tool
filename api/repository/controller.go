package repository

import (
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis" // Importing Go Redis package
	"github.com/google/uuid"
	"github.com/shivesh/URL_Shortner/database"
	"github.com/shivesh/URL_Shortner/helpers"
)

// ResolveURL resolves the short URL and redirects to the original URL.
func (repo *Repository) ResolveURL(c *gin.Context) {
	shortUrl := c.Param("url")

	// Get the original URL from the short URL in the Redis database
	url, err := repo.ShortUrlDBClient.Get(shortUrl).Result()
	if err == redis.Nil {
		c.JSON(404, gin.H{
			"Error": "No URL mapped to this short URL",
		})
		return
	} else if err != nil {
		c.JSON(500, gin.H{
			"Error": "Error while connecting to the database",
		})
		return
	}

	_ = repo.RateLimitDBClient.Incr("counter") // Increment the counter for rate limiting

	c.Redirect(301, url) // Redirect to the original URL
}

// ShortenURL shortens the given URL and provides a custom short URL.
func (repo *Repository) ShortenURL(c *gin.Context) {
	body := database.Request{}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(500, gin.H{
			"Error": err,
		})
		return
	}

	// Rate limiting logic
	_, err := repo.RateLimitDBClient.Get(c.ClientIP()).Result()
	if err == redis.Nil {
		_ = repo.RateLimitDBClient.Set(c.ClientIP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else {
		// Handle rate limit exceeding
		value, _ := repo.RateLimitDBClient.Get(c.ClientIP()).Result()
		valueInt, _ := strconv.Atoi(value)
		if valueInt <= 0 {
			limit, _ := repo.RateLimitDBClient.Get(c.ClientIP()).Result()
			limitInt, _ := strconv.Atoi(limit)
			c.JSON(404, gin.H{
				"error": "Rate limit reached",
				"rate":  limitInt / int(time.Nanosecond) / int(time.Minute),
			})
			return
		}
	}

	// Check if the input URL is a valid URL
	if !govalidator.IsURL(body.Url) {
		c.JSON(400, gin.H{"Error": "Invalid URL"})
	}

	// Handle domain error
	if !helpers.RemoveDomainError(body.Url) {
		c.JSON(400, gin.H{"Error": "Remove domain error"})
	}

	// Enforce HTTPS
	body.Url = helpers.EnforceHTTP(body.Url)

	// Short URL generation logic
	var id string
	if body.Shorturl == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.Shorturl
	}

	// Check if the custom short URL is already in use
	val, _ := repo.ShortUrlDBClient.Get(id).Result()
	if val != "" {
		c.JSON(403, gin.H{"Error": "URL custom short already in use"})
		return
	}

	// Set the shortened URL with an expiry time
	if body.Expiry == 0 {
		body.Expiry = 24
	}
	err = repo.ShortUrlDBClient.Set(id, body.Url, body.Expiry*3600*time.Second).Err()
	if err != nil {
		c.JSON(500, gin.H{
			"Error": "Unable to connect to the database",
		})
		return
	}

	// Prepare the response
	resp := database.Response{
		URL:              body.Url,
		CustomedShortURL: "",
		Expiry:           body.Expiry,
		RateRemaining:    10,
		RateLimitReset:   30,
	}

	// Decrement the rate limit count
	repo.RateLimitDBClient.Decr(c.ClientIP())

	val, _ = repo.RateLimitDBClient.Get(c.ClientIP()).Result()
	resp.RateRemaining, _ = strconv.Atoi(val)

	// Calculate rate limit reset time
	ttl, _ := repo.RateLimitDBClient.TTL(c.ClientIP()).Result()
	resp.RateLimitReset = ttl / time.Nanosecond / time.Minute

	resp.CustomedShortURL = os.Getenv("DOMAIN") + "/" + id

	c.JSON(200, resp)
}
