package main

import (
	"log"
	"os"

	http_interfaces "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	container "github.com/openlabun/CODER/apps/api_v2/internal/application/container"
	)

func newFiberApp(appContainer *container.Application) *fiber.App {
	app := fiber.New()
	_ = appContainer

	allowOrigins := os.Getenv("APIV2_CORS_ALLOW_ORIGINS")
	if allowOrigins == "" {
		allowOrigins = "*"
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins: allowOrigins,
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization,X-User-Email",
	}))

	http_interfaces.RegisterRoutes(app, appContainer)

	return app
}

func main() {
	application, err := container.BuildApplicationContainer()
	if err != nil {
		log.Fatalf("bootstrap failed: %v", err)
	}

	app := newFiberApp(application)

	//endpoint de documentacion

	port := os.Getenv("APIV2_PORT")
	if port == "" {
		port = "4000"
	}

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("server exited with error: %v", err)
	}
}
