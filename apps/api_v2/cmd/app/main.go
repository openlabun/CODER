package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	container "github.com/openlabun/CODER/apps/api_v2/internal/application/container"
	roble_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
	roble_user_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/user"
	rabbitmq_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/publisher/rabbitMQ"
	security_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/security"
	http_interfaces "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http"

	course_repository "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/course"
	exam_repository "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/exam"
	submission_repository "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/submission"
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

	courseRepository := course_repository.NewCourseRepository(robleAdapter)
	examRepository := exam_repository.NewExamRepository(robleAdapter)
	challengeRepository := exam_repository.NewChallengeRepository(robleAdapter)
	testCaseRepository := exam_repository.NewTestCaseRepository(robleAdapter)

	submissionRepository := submission_repository.NewSubmissionRepository(robleAdapter)
	sessionRepository := submission_repository.NewSessionRepository(robleAdapter)
	submissionResRepository := submission_repository.NewSubmissionResultRepository(robleAdapter)

	publisherAdapter, err := rabbitmq_infrastructure.NewRabbitMQAdapter()
	if err != nil {
		return nil, fmt.Errorf("initialize rabbitmq publisher adapter: %w", err)
	}

	deps := container.NewApplicationDependencies(
		authAdapter,
		authAdapter,
		userRepository,
		authAdapter,

		passwordHasher,

		userRepository,
		courseRepository,

		examRepository,
		challengeRepository,
		testCaseRepository,

		submissionRepository,
		sessionRepository,
		submissionResRepository,
		publisherAdapter,
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
