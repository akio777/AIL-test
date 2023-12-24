package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var Validate = validator.New()

func ValidateBody[T any](c *fiber.Ctx) error {
	var body T
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "request body invalid",
		})
	}

	if err := ValidateStruct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Locals("body", body)
	return c.Next()
}

// ValidateStruct is a helper function that calls the global validator instance.
func ValidateStruct(s interface{}) error {
	return Validate.Struct(s)
}
