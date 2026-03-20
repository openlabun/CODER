package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"

	container "github.com/openlabun/CODER/apps/api_v2/internal/application/container"
	ports "github.com/openlabun/CODER/apps/api_v2/internal/application/ports/user"
	roble_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
	roble_user_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/user"
)

var injectedPasswordHasher ports.PasswordHasherPort

// SetPasswordHasher allows external composition roots/tests to inject a hasher
// implementation without coupling this bootstrap to a concrete adapter.
func SetPasswordHasher(passwordHasher ports.PasswordHasherPort) {
	injectedPasswordHasher = passwordHasher
}

func buildApplication(passwordHasher ports.PasswordHasherPort) (*container.Application, error) {
	httpClient := &http.Client{Timeout: 15 * time.Second}

	robleClient, err := roble_infrastructure.NewRobleClient(httpClient)
	if err != nil {
		return nil, fmt.Errorf("initialize roble client: %w", err)
	}

	robleAdapter := roble_infrastructure.NewRobleDatabaseAdapter(robleClient)
	userRepository := roble_user_infrastructure.NewUserRepository(robleAdapter)
	authAdapter := roble_user_infrastructure.NewRobleAuthAdapter(robleAdapter, userRepository)

	deps := container.NewApplicationDependencies(
		authAdapter,
		authAdapter,
		userRepository,
		authAdapter,
		passwordHasher,
	)

	appContainer, err := container.NewApplication(deps)
	if err != nil {
		return nil, fmt.Errorf("initialize application container: %w", err)
	}

	return appContainer, nil
}

func newFiberApp(appContainer *container.Application) *fiber.App {
	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": "api_v2"})
	})

	app.Get("/ready", func(c *fiber.Ctx) error {
		if appContainer == nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "not_ready",
				"reason": "application container is not initialized",
			})
		}

		return c.JSON(fiber.Map{"status": "ready"})
	})

	return app
}

func main() {
	application, err := buildApplication(injectedPasswordHasher)
	if err != nil {
		log.Fatalf("bootstrap failed: %v", err)
	}

	app := newFiberApp(application)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("server exited with error: %v", err)
	}
}
