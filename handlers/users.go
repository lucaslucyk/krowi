package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lucaslucyk/krowi/models"
)

func GetMe(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	return c.Status(fiber.StatusOK).JSON(user)
}
