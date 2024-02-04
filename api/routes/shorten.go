package routes

import (
	"github.com/aaaaayushh/shorten-url-fiber-redis/helpers"
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"time"
)

/*
The json:"" thing in the given code snippet is a struct tag in Golang. It is used to specify the key name used
when marshaling an instance of the struct into JSON. For example, in the request struct, the field URL will be marshaled
into JSON with the key "url" and the field CustomShort will be marshaled with the key "short"

This is a common practice in Golang to provide explicit mapping between struct fields and JSON keys when encoding or
decoding JSON data. It allows developers to have more control over the JSON representation of the struct.
If you're working with JSON in Golang, using struct tags like json:"" can help ensure that the JSON produced by your
code conforms to a specific format or external API requirements.
*/

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int32         `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"x_rate_limit_reset"`
}

func ShortenURL(c *fiber.Ctx) error {

	body := new(request)

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	//implement rate limiting @ 10/30min
	//check if input is an actual URL

	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid URL"})
	}

	//check for domain error
	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Don't be smart :)"})
	}

	//enforce https,SSL
	body.URL = helpers.EnforceHTTP(body.URL)

}
