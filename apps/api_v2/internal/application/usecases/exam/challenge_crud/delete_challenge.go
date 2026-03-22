package exam_usecases

import (
	"context"
	"fmt"

	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	
	user_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	repositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
)

type DeleteChallengeUseCase struct {
	challengeRepository repositories.ChallengeRepository
	userRepository      userRepository.UserRepository
	examRepository      examRepository.ExamRepository
}

func NewDeleteChallengeUseCase(challengeRepository repositories.ChallengeRepository, userRepository userRepository.UserRepository, examRepository examRepository.ExamRepository) *DeleteChallengeUseCase {
	return &DeleteChallengeUseCase{challengeRepository: challengeRepository, userRepository: userRepository, examRepository: examRepository}
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
	
	// [STEP 3] Validate that exam exists and belongs to the teacher
	exam, err := uc.examRepository.GetExamByID(ctx, existingChallenge.ExamID)
	if err != nil {
		return err
	}

	if exam == nil {
		return fmt.Errorf("exam with id %q does not exist", existingChallenge.ExamID)
	}

	if exam.ProfessorID != user.ID {
		return fmt.Errorf("user does not have permissions to delete this challenge")
	}

	// [STEP 4] Delete challenge with user provided values
	err = uc.challengeRepository.DeleteChallenge(ctx, input.ChallengeID)
	if err != nil {
		return err
	}
	
	return nil
}