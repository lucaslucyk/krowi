package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lucaslucyk/krowi/pkg/handlers"
	"github.com/lucaslucyk/krowi/pkg/middleware"
)

func UsersRouter(app *fiber.App, path string) *fiber.Router {
	// create router group
	router := app.Group(path)

	// add endpoints
	router.Get("/me", middleware.DeserializeUser, handlers.GetMe)

	// return pointer to general purpouse
	return &router
}
