package exam_usecases

import (
	"context"
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	mapper "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam/mapper"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"

	user_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/user"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
)

type UpdateExamItemUseCase struct {
	userRepository userRepository.UserRepository
	examRepository examRepository.ExamRepository
	examItemRepository examRepository.ExamItemRepository
	challengeRepository examRepository.ChallengeRepository
}

func NewUpdateExamItemUseCase(userRepository userRepository.UserRepository, examRepository examRepository.ExamRepository, examItemRepository examRepository.ExamItemRepository, challengeRepository examRepository.ChallengeRepository) *UpdateExamItemUseCase {
	return &UpdateExamItemUseCase{
		userRepository: userRepository,
		examRepository: examRepository,
		examItemRepository: examItemRepository,
		challengeRepository: challengeRepository,
	}
}

func (uc *UpdateExamItemUseCase) Execute(ctx context.Context, input dtos.UpdateExamItemInput) (*Entities.ExamItem, error) {
	// [STEP 1] Verify user is teacher and has permissions to create an exam
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

	if user.Role != user_constants.UserRoleProfessor {
		return nil, fmt.Errorf("user does not have permissions to update an exam item")
	}

	// [STEP 2] Validate Exam Item exists
	examItem, err := uc.examItemRepository.GetExamItemByID(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get exam item: %v", err)
	}
	if examItem == nil {
		return nil, fmt.Errorf("exam item with ID %q does not exist", input.ID)
	}

	// [STEP 3] Validate Exam exists and user is owner
	exam, err := uc.examRepository.GetExamByID(ctx, examItem.ExamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get exam: %v", err)
	}

	if exam == nil {
		return nil, fmt.Errorf("exam with ID %q does not exist", examItem.ExamID)
	}

	if exam.ProfessorID != user.ID {
		return nil, fmt.Errorf("user is not the owner of the exam")
	}

	// [STEP 4] Update ExamItem
	updatedExamItem, err := mapper.MapUpdateExamItemInputToExamItemEntity(examItem, input)
	if err != nil {
		return nil, fmt.Errorf("failed to map exam item: %v", err)
	}

	// [STEP 5] Save ExamItem in database
	createdExamItem, err := uc.examItemRepository.UpdateExamItem(ctx, updatedExamItem)
	if err != nil {
		return nil, fmt.Errorf("failed to update exam item: %v", err)
	}


	return createdExamItem, nil
}