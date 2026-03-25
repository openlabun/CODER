package container

import (
	"fmt"
	"net/http"
	"time"

	roble_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
	roble_user_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/user"
	rabbitmq_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/publisher/rabbitMQ"
	security_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/security"
	
	course_repository "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/course"
	exam_repository "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/exam"
	submission_repository "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/submission"

)

func BuildApplicationContainer() (*Application, error) {
	dependencies, err := BuildDependencies()
	if err != nil {
		return nil, fmt.Errorf("failed to build application container: %w", err)
	}

	Application, err := NewApplication(*dependencies)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize application container: %w", err)
	}

	return Application, nil
}

func BuildDependencies() (*ApplicationDependencies, error) {
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
	examItemRepository := exam_repository.NewExamItemRepository(robleAdapter)
	testCaseRepository := exam_repository.NewTestCaseRepository(robleAdapter)

	submissionRepository := submission_repository.NewSubmissionRepository(robleAdapter)
	sessionRepository := submission_repository.NewSessionRepository(robleAdapter)
	submissionResRepository := submission_repository.NewSubmissionResultRepository(robleAdapter)

	publisherAdapter, err := rabbitmq_infrastructure.NewRabbitMQAdapter()
	if err != nil {
		return nil, fmt.Errorf("initialize rabbitmq publisher adapter: %w", err)
	}

	dependencies := NewApplicationDependencies(
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
		examItemRepository,

		submissionRepository,
		sessionRepository,
		submissionResRepository,
		publisherAdapter,
	)

	return &dependencies, nil
}