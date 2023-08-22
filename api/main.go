package main

import (
	"github.com/gin-gonic/gin"                  // Importing the Gin web framework package
	"github.com/shivesh/URL_Shortner/bootstrap" // Importing a custom package named bootstrap
)

func main() {
	// Initialize the Gin web framework with default middleware (logging, recovery)
	app := gin.Default()

	// Initialize the application using the bootstrap package
	bootstrap.InitializeApp(app)
}
