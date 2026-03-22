package submission_usecases

import (
	"context"
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"

	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	submissionRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
	user_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
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

	role := user.Role

	// [STEP 2] Verify if user is a student, only query for his own submissions
	if role == user_entities.UserRoleStudent || input.UserID != user.ID {
		return nil, fmt.Errorf("user does not have permissions to view submissions for this user")
	}

	// [STEP 3] Get all submissions for the user
	submissions, err := uc.submissionRepository.GetSubmissionsByUserID(ctx, user.ID)
		
	if err != nil {
		return nil, err
	}

	

	// [STEP 4]  

	return submissions, nil
}
