package repository

import "github.com/gin-gonic/gin" // Importing the Gin web framework package

// SetupRoutes sets up the routes for the repository using the provided Gin engine.
func (repo *Repository) SetupRoutes(app *gin.Engine) {
	// Define a GET route to resolve shortened URLs
	app.GET("/:url", repo.ResolveURL)

	// Define a POST route to shorten URLs
	app.POST("/shorten", repo.ShortenURL)
}
