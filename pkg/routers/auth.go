package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lucaslucyk/krowi/pkg/handlers"
)

func AuthRouter(app *fiber.App, path string) *fiber.Router {
	// create router group
	router := app.Group(path)

	router.Post("/register", handlers.SignUp)
	router.Post("/login", handlers.LogIn)
	router.Get("/logout", handlers.Logout)

	// return pointer to general purpouse
	return &router
}
