package exam_usecases

import (
	"context"
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"

	user_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/user"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
)

type DeleteExamItemUseCase struct {
	userRepository userRepository.UserRepository
	examRepository examRepository.ExamRepository
	examItemRepository examRepository.ExamItemRepository
	challengeRepository examRepository.ChallengeRepository
}

func NewDeleteExamItemUseCase(userRepository userRepository.UserRepository, examRepository examRepository.ExamRepository, examItemRepository examRepository.ExamItemRepository, challengeRepository examRepository.ChallengeRepository) *DeleteExamItemUseCase {
	return &DeleteExamItemUseCase{
		userRepository: userRepository,
		examRepository: examRepository,
		examItemRepository: examItemRepository,
		challengeRepository: challengeRepository,
	}
}

func (uc *DeleteExamItemUseCase) Execute(ctx context.Context, input dtos.DeleteExamItemInput) error {
	// [STEP 1] Verify user is teacher and has permissions to create an exam
	userEmail, err := services.UserEmailFromContext(ctx)
	if err != nil {
		return err
	}

	user, err := uc.userRepository.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return err
	}

	if user == nil {
		return fmt.Errorf("user with email %q does not exist", userEmail)
	}

	if user.Role != user_constants.UserRoleProfessor {
		return fmt.Errorf("user does not have permissions to delete an exam item")
	}

	// [STEP 2] Validate Exam Item exists
	examItem, err := uc.examItemRepository.GetExamItemByID(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("failed to get exam item: %v", err)
	}
	if examItem == nil {
		return fmt.Errorf("exam item with ID %q does not exist", input.ID)
	}

	// [STEP 3] Validate Exam exists and user is owner
	exam, err := uc.examRepository.GetExamByID(ctx, examItem.ExamID)
	if err != nil {
		return fmt.Errorf("failed to get exam: %v", err)
	}
	if exam == nil {
		return fmt.Errorf("exam with ID %q does not exist", examItem.ExamID)
	}
	if exam.ProfessorID != user.ID {
		return fmt.Errorf("user is not the owner of the exam")
	}

	// [STEP 4] Delete ExamItem in database
	err = uc.examItemRepository.DeleteExamItem(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("failed to delete exam item: %v", err)
	}

	return nil
}