package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/lucaslucyk/krowi/database"
	router "github.com/lucaslucyk/krowi/routers"
)

func main() {
	// connect to db
	database.Connect()

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
	router.SetupRoutes(app)

	// run server
	log.Fatal(app.Listen(":8000"))
}
