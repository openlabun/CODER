package submission_usecases

import (
	"context"
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
	mapper "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission/mapper"

	user_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	resultsRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
	submissionRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
)

type GetSubmissionStatusUseCase struct {
	userRepository userRepository.UserRepository
	resultsRepository resultsRepository.SubmissionResultRepository
	submissionRepository submissionRepository.SubmissionRepository
}

func NewGetSubmissionStatusUseCase(userRepository userRepository.UserRepository, resultsRepository resultsRepository.SubmissionResultRepository, submissionRepository submissionRepository.SubmissionRepository) *GetSubmissionStatusUseCase {
	return &GetSubmissionStatusUseCase{
		userRepository: userRepository,
		resultsRepository: resultsRepository,
		submissionRepository: submissionRepository,
	}
}

func (uc *GetSubmissionStatusUseCase) Execute(ctx context.Context, input dtos.GetSubmissionStatusInput) (*dtos.SubmissionOutputDTO, error) {
	// [STEP 1] Verify user
	userEmail, err := services.UserEmailFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepository.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("user with email %q does not exist", userEmail)
	}

	if user.Role != user_entities.UserRoleStudent && user.Role != user_entities.UserRoleProfessor && user.Role != user_entities.UserRoleAdmin {
		return nil, fmt.Errorf("user role %q is not allowed to view submission status", user.Role)
	}

	// [STEP 2] Verify that submission exists
	submission, err := uc.submissionRepository.GetSubmissionByID(ctx, input.SubmissionID)
	if err != nil {
		return nil, err
	}

	if submission == nil {
		return nil, fmt.Errorf("submission with id %q does not exist", input.SubmissionID)
	}

	// [STEP 3] If user is a student, only query for his own submission status
	if user.Role == user_entities.UserRoleStudent && submission.UserID != user.ID {
		return nil, fmt.Errorf("user does not have permissions to view submission status for this submission")
	}

	// [STEP 4] Map submission status entity to output DTO and return it
	results, err := uc.getSubmissionResults(ctx, submission)
	if err != nil {
		return nil, err
	}

	dto := mapper.MapSubmissionOutputDTO(submission, results)
	if dto == nil {
		return nil, fmt.Errorf("failed to map submission output DTO")
	}

	return dto, nil
}

func (uc *GetSubmissionStatusUseCase) getSubmissionResults (ctx context.Context, submission *Entities.Submission) ([]Entities.SubmissionResult, error) {
	results, err := uc.resultsRepository.GetResultsBySubmissionID(ctx, submission.ID)
	if err != nil {
		return nil, err
	}

	derefResults := make([]Entities.SubmissionResult, len(results))
	for i, result := range results {
		derefResults[i] = *result
	}

	return derefResults, err
}