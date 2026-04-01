package submission_usecases

import (
	"context"
	"fmt"
	"os"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"

	passwordHasher "github.com/openlabun/CODER/apps/api_v2/internal/application/ports/user"
	userPort "github.com/openlabun/CODER/apps/api_v2/internal/application/ports/user"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	submissionRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	domain_services "github.com/openlabun/CODER/apps/api_v2/internal/domain/services"

	mapper "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission/mapper"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
)

type UpdateResultUseCase struct {
	userService          userPort.LoginPort
	userRepository       userRepository.UserRepository
	submissionRepository submissionRepository.SubmissionRepository
	sessionRepository    submissionRepository.SessionRepository
	challengeRepository  examRepository.ChallengeRepository
	testCaseRepository   examRepository.TestCaseRepository
	ioVariableRepository examRepository.IOVariableRepository
	resultRepository     submissionRepository.SubmissionResultRepository
	passwordHasher       passwordHasher.PasswordHasherPort
}

func NewUpdateResultUseCase(userRepository userRepository.UserRepository, submissionRepository submissionRepository.SubmissionRepository, sessionRepository submissionRepository.SessionRepository, challengeRepository examRepository.ChallengeRepository, testCaseRepository examRepository.TestCaseRepository, ioVariableRepository examRepository.IOVariableRepository, resultRepository submissionRepository.SubmissionResultRepository, userService userPort.LoginPort, passwordHasher passwordHasher.PasswordHasherPort) *UpdateResultUseCase {
	return &UpdateResultUseCase{
		userRepository:       userRepository,
		submissionRepository: submissionRepository,
		sessionRepository:    sessionRepository,
		challengeRepository:  challengeRepository,
		testCaseRepository:   testCaseRepository,
		ioVariableRepository: ioVariableRepository,
		resultRepository:     resultRepository,
		userService:          userService,
		passwordHasher:       passwordHasher,
	}
}

func (uc *UpdateResultUseCase) Execute(ctx context.Context, input dtos.UpdateResultInput) (*Entities.SubmissionResult, error) {
	// [STEP 1] Validate internal service key
	workerKey, ok := services.AccessInternalServiceKeyFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no internal service key provided in context")
	}

	validWorkerKey := os.Getenv("WORKER_KEY")
	if validWorkerKey == "" {
		return nil, fmt.Errorf("internal service key is not configured in environment variables")
	}

	if workerKey != validWorkerKey {
		return nil, fmt.Errorf("invalid internal service key")
	}

	// [STEP 2] Login with internal admin user and set access token in context
	internalEmail := os.Getenv("INTERNAL_USER_EMAIL")
	internalPassword, err := uc.passwordHasher.Hash(os.Getenv("INTERNAL_USER_PASSWORD"))
	if err != nil {
		return nil, fmt.Errorf("failed to hash internal user password: %w", err)
	}

	user, err := uc.userService.LoginUser(internalEmail, internalPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to login with internal user: %w", err)
	}

	ctx = services.WithAccessToken(ctx, user.Token.AccessToken)

	// [STEP 3] Retrieve the submission result from the repository
	submissionResult, err := uc.resultRepository.GetResultByID(ctx, input.ResultID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve submission result: %w", err)
	}

	if submissionResult == nil {
		return nil, fmt.Errorf("submission result with id %q does not exist", input.ResultID)
	}

	// [STEP 4] Get submission test case
	testCase, err := uc.testCaseRepository.GetTestCaseByID(ctx, submissionResult.TestCaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve test case: %w", err)
	}

	// [STEP 4] Update the submission result entity with the new data
	updatedResult, err := mapper.MapResultInputToSubmissionResultEntity(input, submissionResult, testCase)
	if err != nil {
		return nil, fmt.Errorf("failed to map input to submission result entity: %w", err)
	}

	// [STEP 5] Check results if status is executed, it changes status to accepted or wrong_answer
	if updatedResult.Status == Entities.SubmissionStatusExecuted {
		updatedResult, err = services.CheckSubmissionResult(updatedResult, testCase)
		if err != nil {
			return nil, fmt.Errorf("failed to check submission result: %w", err)
		}
	}

	// [STEP 6] Save the updated submission result back to the repository
	result, err := domain_services.UpdateSubmissionResult(ctx, updatedResult, uc.resultRepository, uc.ioVariableRepository)
	if err != nil {
		return nil, fmt.Errorf("failed to update submission result: %w", err)
	}

	return result, nil
}
