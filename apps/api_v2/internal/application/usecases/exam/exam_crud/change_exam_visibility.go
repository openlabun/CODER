package exam_usecases

import (
	"context"
	"fmt"
	
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	user_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
	mapper "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam/mapper"
)

type ChangeExamVisibilityUseCase struct {
	userRepository userRepository.UserRepository
	examRepository examRepository.ExamRepository
}

func NewChangeExamVisibilityUseCase(userRepository userRepository.UserRepository, examRepository examRepository.ExamRepository) *ChangeExamVisibilityUseCase {
	return &ChangeExamVisibilityUseCase{
		userRepository: userRepository,
		examRepository: examRepository,
	}
}

func (uc *ChangeExamVisibilityUseCase) Execute(ctx context.Context, input dtos.ChangeExamVisibilityInput) (*Entities.Exam, error) {
	// [STEP 1] Verify user is teacher and has permissions to change exam visibility
	userEmail, err := services.UserEmailFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepository.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return nil, err
	}

	if user.Role != user_entities.UserRoleProfessor {
		return nil, fmt.Errorf("user does not have permissions to create an exam")
	}

	// [STEP 2] Get exam entity to be updated
	exam, err := uc.examRepository.GetExamByID(ctx, input.ExamID)
	if err != nil {
		return nil, err
	}

	// [STEP 3] Update exam entity with new visibility value
	exam, err = mapper.MapExamVisibilityInputToExamEntity(exam, input)
	if err != nil {
		return nil, err
	}

	// [STEP 4] Save updated exam entity
	exam, err = uc.examRepository.UpdateExam(ctx, exam)
	if err != nil {
		return nil, err
	}

	return exam, nil
}
