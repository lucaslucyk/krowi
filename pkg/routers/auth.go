package router

import (
	"github.com/gofiber/fiber/v2"
	authenticator "github.com/lucaslucyk/krowi/pkg/authenticators"
	"github.com/lucaslucyk/krowi/pkg/handlers"
	"github.com/lucaslucyk/krowi/pkg/middleware"
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

func OAuthRouter(
	app *fiber.App,
	auth *authenticator.Authenticator,
	path string,
) *fiber.Router {
	// create router group
	router := app.Group(path)

	router.Get("/login", handlers.OAuthHandler(auth))
	router.Get("/logout", handlers.OLogoutHandler(auth))
	router.Get("/callback", handlers.CallbackHandler(auth))
	router.Get("/me", handlers.GetOauthMe)
	router.Get("/me2", middleware.EnsureValidToken, handlers.GetOauthMe)

	// return pointer to general purpouse
	return &router
}
