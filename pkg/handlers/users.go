package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/lucaslucyk/krowi/pkg/models"
)

func GetMe(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	return c.Status(fiber.StatusOK).JSON(user)
}

func GetOauthMe(c *fiber.Ctx) error {
	session := c.Locals("session").(*session.Session)
	profile := session.Get("profile")
	return c.Status(fiber.StatusOK).JSON(profile)
}
