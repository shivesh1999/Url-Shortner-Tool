package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/shivesh/URL_Shortner/routes"
)

// setup two routes, one for shortening the url
// the other for resolving the url
func setupRoutes(app *fiber.App) {
	app.Get("/:url", routes.ResolveURL)
	app.Post("/api/v1", routes.ShortenURL)
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
	//Here we use fiber as a web framework which works
	//the same as express in javascript
	app := fiber.New()
	app.Use(logger.New())
	setupRoutes(app)

	//start the server for the reuqired port
	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}
