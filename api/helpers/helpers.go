package helpers

import (
	"os"
	"strings"
)

// EnforceHTTP enforces the use of "http://" prefix for a URL if not present.
func EnforceHTTP(url string) string {
	if url[:4] != "http" {
		return "http://" + url
	}
	return url
}

// RemoveDomainError removes common URL components and checks for domain errors.
func RemoveDomainError(url string) bool {
	// Check if the URL is equal to the specified domain in the environment variable
	if url == os.Getenv("DOMAIN") {
		return false
	}

	// Remove "http://" and "https://" prefixes, and remove "www." prefix if present
	newURL := strings.Replace(url, "http://", "", 1)
	newURL = strings.Replace(newURL, "https://", "", 1)
	newURL = strings.Replace(newURL, "www.", "", 1)

	// Split the URL to get the domain part
	newURL = strings.Split(newURL, "/")[0]

	// Check if the new URL is equal to the specified domain in the environment variable
	if newURL == os.Getenv("DOMAIN") {
		return false
	}

	return true
}
