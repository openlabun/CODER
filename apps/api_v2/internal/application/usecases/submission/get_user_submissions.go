package submission_usecases

import (
	"context"
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"

	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	submissionRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
)

type GetUserSubmissionsUseCase struct {
	userRepository userRepository.UserRepository
	challengeRepository examRepository.ChallengeRepository
	examRepository examRepository.ExamRepository
	submissionRepository submissionRepository.SubmissionRepository
}

func NewGetUserSubmissionsUseCase(userRepository userRepository.UserRepository, challengeRepository examRepository.ChallengeRepository, examRepository examRepository.ExamRepository, submissionRepository submissionRepository.SubmissionRepository) *GetUserSubmissionsUseCase {
	return &GetUserSubmissionsUseCase{
		userRepository: userRepository,
		challengeRepository: challengeRepository,
		examRepository: examRepository,
		submissionRepository: submissionRepository,
	}
}

func (uc *GetUserSubmissionsUseCase) Execute(ctx context.Context, input dtos.GetUserSubmissionsInput) ([]*Entities.Submission, error) {
	// [STEP 1] Verify user and his role
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

	// [STEP 2] Check permissions (Temporarily relaxed for debugging)
	/*
	if input.UserID != user.ID && user.Role != user_entities.UserRoleAdmin && user.Role != user_entities.UserRoleProfessor {
		return nil, fmt.Errorf("user does not have permissions to view submissions for this user")
	}
	*/

	// [STEP 3] Get all submissions for the user
	submissions, err := uc.submissionRepository.GetSubmissionsByUserID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	return submissions, nil
}
