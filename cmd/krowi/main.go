package main

import (
	"encoding/gob"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	authenticator "github.com/lucaslucyk/krowi/pkg/authenticators"
	"github.com/lucaslucyk/krowi/pkg/config"
	"github.com/lucaslucyk/krowi/pkg/database"
	"github.com/lucaslucyk/krowi/pkg/handlers"
	"github.com/lucaslucyk/krowi/pkg/middleware"
	routers "github.com/lucaslucyk/krowi/pkg/routers"
)

func main() {
	// connect to db
	database.Connect()
	cfg, err := config.New()
	if err != nil {
		panic(err.Error())
	}
	// create app
	app := fiber.New()
	gob.Register(map[string]interface{}{})
	store := session.New()
	middleware.SetupSessionStoreMiddleware(app, store)
	middleware.IsOAuthenticatedMiddleware(app)

	// config cors
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173, https://rm1qc2pg-8000.brs.devtunnels.ms",
		AllowMethods:     "GET, POST",
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept",
	}))

	auth, err := authenticator.New()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	app.Get("/", handlers.Home)

	// setup routes
	_ = routers.AuthRouter(app, "/auth")
	_ = routers.OAuthRouter(app, auth, "/oauth")
	_ = routers.UsersRouter(app, "/users")

	// run server
	log.Fatal(app.Listen(fmt.Sprintf(
		"%s:%s",
		cfg.KrowiHost,
		cfg.KrowiPort,
	)))
}
