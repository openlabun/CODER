package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"

	container "github.com/openlabun/CODER/apps/api_v2/internal/application/container"
	roble_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
	roble_user_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/user"
	security_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/security"
	http_interfaces "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http"
)

func buildApplication() (*container.Application, error) {
	// Start clients
	httpClient := &http.Client{Timeout: 15 * time.Second}
	robleClient, err := roble_infrastructure.NewRobleClient(httpClient)
	if err != nil {
		return nil, fmt.Errorf("initialize roble client: %w", err)
	}

	// Start adapters and repositories
	robleAdapter := roble_infrastructure.NewRobleDatabaseAdapter(robleClient)
	userRepository := roble_user_infrastructure.NewUserRepository(robleAdapter)
	authAdapter := roble_user_infrastructure.NewRobleAuthAdapter(robleAdapter, userRepository)
	passwordHasher := security_infrastructure.NewSecurityAdapter()

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
	_ = appContainer

	http_interfaces.RegisterRoutes(app)

	return app
}

func main() {
	application, err := buildApplication()
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
