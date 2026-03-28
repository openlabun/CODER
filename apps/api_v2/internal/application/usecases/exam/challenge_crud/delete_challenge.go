package exam_usecases

import (
	"context"
	"fmt"

	submissionRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	domain_services "github.com/openlabun/CODER/apps/api_v2/internal/domain/services"
	user_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	repositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
)

type DeleteChallengeUseCase struct {
	challengeRepository repositories.ChallengeRepository
	userRepository      userRepository.UserRepository
	examRepository      examRepository.ExamRepository
	testCaseRepository 	examRepository.TestCaseRepository
	examItemRepository 	examRepository.ExamItemRepository
	submissionRepository submissionRepository.SubmissionRepository
	resultsRepository submissionRepository.SubmissionResultRepository
}

func NewDeleteChallengeUseCase(challengeRepository repositories.ChallengeRepository, userRepository userRepository.UserRepository, examRepository examRepository.ExamRepository, testCaseRepository examRepository.TestCaseRepository, examItemRepository examRepository.ExamItemRepository, submissionRepository submissionRepository.SubmissionRepository, resultsRepository submissionRepository.SubmissionResultRepository) *DeleteChallengeUseCase {
	return &DeleteChallengeUseCase{challengeRepository: challengeRepository, userRepository: userRepository, examRepository: examRepository, testCaseRepository: testCaseRepository, examItemRepository: examItemRepository, submissionRepository: submissionRepository, resultsRepository: resultsRepository}
}

func (uc *DeleteChallengeUseCase) Execute(ctx context.Context, input dtos.DeleteChallengeInput) error {
	// [STEP 1] Verify user is teacher and has permissions to delete an exam
	userEmail, err := services.UserEmailFromContext(ctx)
	if err != nil {
		return err
	}

	user, err := uc.userRepository.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return err
	}

	if user.Role != user_entities.UserRoleProfessor {
		return fmt.Errorf("user does not have permissions to create an exam")
	}

	// [STEP 2] Validate that the challenge exists
	existingChallenge, err := uc.challengeRepository.GetChallengeByID(ctx, input.ChallengeID)
	if err != nil {
		return err
	}

	if existingChallenge == nil {
		return fmt.Errorf("challenge with id %s does not exist", input.ChallengeID)
	}
	
	// [STEP 3] Validate that challenge belongs to teacher
	if existingChallenge.UserID != user.ID {
		return fmt.Errorf("user does not have permissions to delete this challenge")
	}

	// [STEP 4] Delete challenge with user provided values
	err = domain_services.RemoveChallenge(ctx, input.ChallengeID, uc.challengeRepository, uc.testCaseRepository, uc.examItemRepository, uc.submissionRepository, uc.resultsRepository)
	if err != nil {
		return fmt.Errorf("failed to delete challenge with id %q: %v", input.ChallengeID, err)
	}

	return nil
}