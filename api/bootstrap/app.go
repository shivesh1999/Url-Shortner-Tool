package bootstrap

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"                   // Importing the Gin web framework package
	"github.com/joho/godotenv"                   // Importing the godotenv package for handling environment variables
	"github.com/shivesh/URL_Shortner/database"   // Importing a custom package named database
	"github.com/shivesh/URL_Shortner/repository" // Importing a custom package named repository
)

// InitializeApp initializes the application by setting up routes and starting the Gin engine.
func InitializeApp(app *gin.Engine) {
	// Check if the APP_ENV environment variable is set
	_, ok := os.LookupEnv("APP_ENV")

	// If APP_ENV is not set, load environment variables from .env file
	if !ok {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal(err)
		}
	}

	// Create separate database clients for rate limiting and short URLs
	rateLimitClient := database.CreateClient(1)
	shortUrlClient := database.CreateClient(0)

	// Create a repository instance with the database clients
	repo := repository.Repository{
		RateLimitDBClient: rateLimitClient,
		ShortUrlDBClient:  shortUrlClient,
	}

	// Setup routes using the repository and the provided Gin engine
	repo.SetupRoutes(app)

	// Start the Gin engine and listen for incoming requests
	log.Fatal(app.Run(os.Getenv("APP_PORT")))
}
