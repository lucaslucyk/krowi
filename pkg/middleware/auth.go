package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/lucaslucyk/krowi/pkg/config"
	"github.com/lucaslucyk/krowi/pkg/database"
	"github.com/lucaslucyk/krowi/pkg/models"
)

func DeserializeUser(c *fiber.Ctx) error {
	var (
		token string
		cfg   *config.Config
		err   error
	)

	cfg, err = config.New()
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{"message": "Not available now"})
	}

	authorization := c.Get("Authorization")

	if strings.HasPrefix(authorization, "Bearer ") {
		token = strings.TrimPrefix(authorization, "Bearer ")
	} else if c.Cookies("token") != "" {
		token = c.Cookies("token")
	}

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{"message": "Not authenticated"})
	}

	tokenByte, err := jwt.Parse(
		token,
		func(jwtToken *jwt.Token) (interface{}, error) {
			if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("sign error: %s", jwtToken.Header["alg"])
			}

			return []byte(cfg.SecretKey), nil
		})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"message": fmt.Sprintf("invalid token: %v", err),
			})
	}

	claims, ok := tokenByte.Claims.(jwt.MapClaims)
	if !ok || !tokenByte.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"message": "invalid token claim",
			})
	}

	var user models.User
	db := database.DB
	db.First(&user, "id = ?", fmt.Sprint(claims["sub"]))

	if user.ID.String() != claims["sub"] {
		return c.Status(fiber.StatusForbidden).JSON(
			fiber.Map{
				"message": "user not found",
			})
	}

	c.Locals("user", &user)

	return c.Next()
}
