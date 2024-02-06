package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/lucaslucyk/krowi/pkg/config"
	"github.com/lucaslucyk/krowi/pkg/database"
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

	// config cors
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowMethods:     "GET, POST",
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept",
	}))

	// setup routes
	_ = routers.AuthRouter(app, "/auth")
	_ = routers.UsersRouter(app, "/users")

	// run server
	log.Fatal(app.Listen(fmt.Sprintf(
		"%s:%s",
		cfg.KrowiHost,
		cfg.KrowiPort,
	)))
}
