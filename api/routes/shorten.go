package routes

import (
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shivesh/URL_Shortner/database"
	"github.com/shivesh/URL_Shortner/helpers"
)

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

type response struct {
	URL            string        `json:"url"`
	CustomedShort  string        `json:"short"`
	Expiry         time.Duration `json:"expiry"`
	RateRemaining  int           `json:"rate_limit"`
	RateLimitReset time.Duration `json:"rate_limit_reset"`
}

func ShortenURL(c *fiber.Ctx) error {

	body := new(request)
	if err := c.BodyParser(&body); err != nil {
		c.Status(fiber.StatusBadRequest).JSON("error : unable to parse json")
	}

	//implement rate limit

	r2 := database.CreateClient(1)
	defer r2.Close()
	_, err := r2.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil {
		_ = r2.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else {
		value, _ := r2.Get(database.Ctx, c.IP()).Result()
		valueInt, _ := strconv.Atoi(value)
		if valueInt <= 0 {
			limit, _ := r2.Get(database.Ctx, c.IP()).Result()
			limitInt, _ := strconv.Atoi(limit)
			return c.Status(fiber.StatusServiceUnavailable).JSON(&fiber.Map{
				"error": "Rate limit reached",
				"rate":  limitInt / int(time.Nanosecond) / int(time.Minute),
				"value": valueInt,
			})
		}
	}

	//check for the input if a actual url

	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON("error : invalid URL")
	}

	//chek for domain error

	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON("error : remove domain error")
	}

	//enforce https, SSL

	body.URL = helpers.EnforceHTTP(body.URL)

	var id string
	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}
	r := database.CreateClient(0)
	defer r.Close()

	val, _ := r.Get(database.Ctx, id).Result()
	if val != "" {
		return c.Status(fiber.StatusForbidden).JSON(&fiber.Map{
			"error": "URL custom shorl already in use",
		})
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}

	err = r.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "unable to connect to the docker file",
		})
	}

	resp := response{
		URL:            body.URL,
		CustomedShort:  "",
		Expiry:         body.Expiry,
		RateRemaining:  10,
		RateLimitReset: 30,
	}

	r2.Decr(database.Ctx, c.IP())

	val, _ = r2.Get(database.Ctx, c.IP()).Result()
	resp.RateRemaining, _ = strconv.Atoi(val)

	ttl, _ := r2.TTL(database.Ctx, c.IP()).Result()
	resp.RateLimitReset = ttl / time.Nanosecond / time.Minute

	resp.CustomedShort = os.Getenv("DOMAIN") + "/" + id

	return c.Status(fiber.StatusOK).JSON(resp)
}
