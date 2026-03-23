package submission_usecases

import (
	"context"
	"fmt"
	"os"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	submissionRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"

	mapper "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission/mapper"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
)

type UpdateResultUseCase struct {
	userRepository userRepository.UserRepository
	submissionRepository submissionRepository.SubmissionRepository
	sessionRepository submissionRepository.SessionRepository
	challengeRepository examRepository.ChallengeRepository
	testCaseRepository examRepository.TestCaseRepository
	resultRepository submissionRepository.SubmissionResultRepository
}

func NewUpdateResultUseCase(userRepository userRepository.UserRepository, submissionRepository submissionRepository.SubmissionRepository, sessionRepository submissionRepository.SessionRepository, challengeRepository examRepository.ChallengeRepository, testCaseRepository examRepository.TestCaseRepository, resultRepository submissionRepository.SubmissionResultRepository) *UpdateResultUseCase {
	return &UpdateResultUseCase{
		userRepository: userRepository,
		submissionRepository: submissionRepository,
		sessionRepository: sessionRepository,
		challengeRepository: challengeRepository,
		testCaseRepository: testCaseRepository,
		resultRepository: resultRepository,
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

	// [STEP 1] Retrieve the submission result from the repository
	submissionResult, err := uc.resultRepository.GetResultByID(ctx, input.ResultID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve submission result: %w", err)
	}

	// [STEP 2] Update the submission result entity with the new data
	updatedResult, err := mapper.MapResultInputToSubmissionResultEntity(input, submissionResult)
	if err != nil {
		return nil, fmt.Errorf("failed to map input to submission result entity: %w", err)
	}

	// [STEP 3] Save the updated submission result back to the repository
	result, err := uc.resultRepository.UpdateResult(ctx, updatedResult)
	if err != nil {
		return nil, fmt.Errorf("failed to update submission result: %w", err)
	}

	return result, nil
}

