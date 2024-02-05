package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lucaslucyk/krowi/pkg/handlers"
	"github.com/lucaslucyk/krowi/pkg/middleware"
)

func SetupRoutes(app *fiber.App) {
	app.Post("/register", handlers.SignUp)
	app.Post("/login", handlers.LogIn)
	app.Get("/logout", handlers.Logout)
	app.Get("/me", middleware.DeserializeUser, handlers.GetMe)
}
