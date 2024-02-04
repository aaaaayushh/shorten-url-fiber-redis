package routes

import (
	"github.com/aaaaayushh/shorten-url-fiber-redis/database"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"log"
)

func ResolveURL(c *fiber.Ctx) error {
	url := c.Params("url")

	r := database.CreateClient(0)
	defer func(r *redis.Client) {
		err := r.Close()
		if err != nil {
			log.Fatal("error closing db connection")
		}
	}(r) //close db connection before function exits

	value, err := r.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "short not found"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot connect to DB"})
	}

	rInr := database.CreateClient(1)
	_ = rInr.Incr(database.Ctx, "counter")

	return c.Redirect(value, 301)

}
