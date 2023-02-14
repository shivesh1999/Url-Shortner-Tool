package routes

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/shivesh/URL_Shortner/database"
)

// function tp resolve the shoetened url and redirect to the main url
func ResolveURL(c *fiber.Ctx) error {
	url := c.Params("url")

	r := database.CreateClient(0)
	defer r.Close()

	//search for the url in the database in dbno 0
	//dbno 0 is a redis key pair database wherre key is the short url
	//and values is the main url to which we want to redirect
	value, err := r.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"error": "short not found in the database",
		})
	} else if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"error ": "cannot connect to database",
		})
	}

	rInr := database.CreateClient(1)
	defer rInr.Close()
	_ = rInr.Incr(database.Ctx, "counter")

	return c.Redirect(value, 301)
}
