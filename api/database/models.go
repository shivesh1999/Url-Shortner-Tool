package database

import "time"

// Request represents the structure for incoming URL shortening requests.
type Request struct {
	Url      string        `json:"url"`      // Original URL to be shortened
	Shorturl string        `json:"shorturl"` // Custom short URL (optional)
	Expiry   time.Duration `json:"expiry"`   // Expiry duration for the short URL (in hours)
}

// Response represents the structure for URL shortening response.
type Response struct {
	URL              string        `json:"url"`              // Original URL
	CustomedShortURL string        `json:"shortURL"`         // Shortened URL with custom identifier
	Expiry           time.Duration `json:"expiry"`           // Expiry duration for the short URL (in hours)
	RateRemaining    int           `json:"rate_limit"`       // Remaining rate limit for the client
	RateLimitReset   time.Duration `json:"rate_limit_reset"` // Time until rate limit resets
}
