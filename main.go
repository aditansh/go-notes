package main

import (
	"fmt"
	"log"

	initializers "github.com/aditansh/go-notes/config"
	database "github.com/aditansh/go-notes/db"
	"github.com/aditansh/go-notes/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	// err := sentry.Init(sentry.ClientOptions{
	// 	Dsn: os.Getenv("SENTRY_DSN"),
	// 	// Set TracesSampleRate to 1.0 to capture 100%
	// 	// of transactions for performance monitoring.
	// 	// We recommend adjusting this value in production,
	// 	TracesSampleRate: 1.0,
	// })
	// if err != nil {
	// 	log.Fatalf("sentry.Init: %s", err)
	// }

	app := fiber.New()

	config, err := initializers.LoadEnvVariables(".")
	if err != nil {
		log.Fatalf("Error loading env variables: %s", err.Error())
	}
	database.ConnectDB(&config)

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.ClientOrigin,
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST, PUT, DELETE",
		AllowCredentials: true,
	}))

	fmt.Println("Server started")

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  true,
			"message": "pong",
		})
	})

	routes.SetupRoutes(app)

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"status":  false,
			"message": "Route Not Found",
		})
	})

	log.Fatal(app.Listen(config.Port))
}
