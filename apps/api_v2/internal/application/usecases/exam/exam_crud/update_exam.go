package exam_usecases

import (
	"context"
	"fmt"
	
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	user_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	exam_validations "github.com/openlabun/CODER/apps/api_v2/internal/domain/validations/exam"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
	mapper "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam/mapper"
)

type UpdateExamUseCase struct {
	userRepository userRepository.UserRepository
	examRepository examRepository.ExamRepository
}

func NewUpdateExamUseCase(userRepository userRepository.UserRepository, examRepository examRepository.ExamRepository) *UpdateExamUseCase {
	return &UpdateExamUseCase{
		userRepository: userRepository,
		examRepository: examRepository,
	}
}

func (uc *UpdateExamUseCase) Execute(ctx context.Context, input dtos.UpdateExamInput) (*Entities.Exam, error) {
	// [STEP 1] Verify user is teacher and has permissions to update an exam
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

	if user.Role != user_entities.UserRoleProfessor {
		return nil, fmt.Errorf("user does not have permissions to update an exam")
	}

	// [STEP 2] Get original exam entity with user provided values
	exam, err := uc.examRepository.GetExamByID(ctx, input.ExamID)
	if err != nil {
		return nil, err
	}

	if exam == nil {
		return nil, fmt.Errorf("exam with id %q does not exist", input.ExamID)
	}

	// [STEP 3] Update exam entity with user provided values
	exam, err = mapper.MapUpdateExamInputToExamEntity(exam, input)
	if err != nil {
		return nil, err
	}

	// [STEP 4] If EndTime is being updated, validate it is not in the past
	if input.EndTime != nil {
		now := services.Now()
		if err := exam_validations.ValidateExamEndTime(exam, now); err != nil {
			return nil, err
		}
	}

	// [STEP 5] Save updated exam entity
	exam, err = uc.examRepository.UpdateExam(ctx, exam)
	if err != nil {
		return nil, err
	}

	return exam, nil
}
