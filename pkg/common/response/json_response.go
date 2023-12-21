package common

import "github.com/gofiber/fiber/v2"

func JSONResponse(c *fiber.Ctx, data interface{}, status int) error {
	response := map[string]interface{}{"data": data}
	c.Status(status)
	return c.JSON(response)
}
func JSONResponseError(c *fiber.Ctx, data interface{}, status int) error {
	response := map[string]interface{}{"error": data}
	c.Status(status)
	return c.JSON(response)
}
