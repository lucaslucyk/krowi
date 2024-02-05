package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/lucaslucyk/krowi/pkg/config"
	"github.com/lucaslucyk/krowi/pkg/models"
	users "github.com/lucaslucyk/krowi/pkg/services"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *fiber.Ctx) error {
	var payload *models.SignUpInput

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"status": "fail", "message": err.Error()})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(payload.Password), bcrypt.DefaultCost)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"status": "fail", "message": err.Error()})
	}

	newUser := models.User{
		Name:     payload.Name,
		Email:    strings.ToLower(payload.Email),
		Password: string(hashedPassword),
	}

	err = users.CreateUser(&newUser)
	if err != nil {
		if err.Error() == users.ALREADY_EXISTS {
			return c.Status(fiber.StatusConflict).JSON(
				fiber.Map{"status": "fail", "message": "User already exists"})
		}

		return c.Status(fiber.StatusBadGateway).JSON(
			fiber.Map{"status": "error", "message": "Something was wrong"})
	}

	return c.Status(fiber.StatusCreated).JSON(&newUser)
}

func LogIn(c *fiber.Ctx) error {
	var err error
	cfg, err := config.New()
	if err != nil {
		panic(err.Error())
	}

	var payload *models.SignInInput
	if err = c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"status": "fail", "message": err.Error()})
	}

	var user models.User
	err = users.GetUserByEmail(&user, strings.ToLower(payload.Email))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"message": "Invalid email or password",
			})
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"message": "Invalid email or password",
			})
	}

	tokenByte := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()
	claims := tokenByte.Claims.(jwt.MapClaims)
	expDuration := time.Hour * 24

	claims["sub"] = user.ID
	claims["exp"] = now.Add(expDuration).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := tokenByte.SignedString([]byte(cfg.SecretKey))

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(
			fiber.Map{
				"status":  "fail",
				"message": fmt.Sprintf("JWT failed: %v", err),
			})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		MaxAge:   60 * 60,
		Secure:   false,
		HTTPOnly: true,
		Domain:   "localhost",
	})

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"token_type":   "bearer",
			"expires_in":   expDuration.Seconds(),
			"access_token": token,
		})
}

func Logout(c *fiber.Ctx) error {
	expired := time.Now().Add(-time.Hour * 24)
	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Value:   "",
		Expires: expired,
	})
	return c.Status(fiber.StatusOK).SendString("")
}
