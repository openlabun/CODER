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

type CloseExamUseCase struct {
	userRepository userRepository.UserRepository
	examRepository examRepository.ExamRepository
}

func NewCloseExamUseCase(userRepository userRepository.UserRepository, examRepository examRepository.ExamRepository) *CloseExamUseCase {
	return &CloseExamUseCase{
		userRepository: userRepository,
		examRepository: examRepository,
	}
}

func (uc *CloseExamUseCase) Execute(ctx context.Context, input dtos.CloseExamInput) (*Entities.Exam, error) {
	// [STEP 1] Verify user is teacher and has permissions to close an exam
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
		return nil, fmt.Errorf("user does not have permissions to close an exam")
	}

	// [STEP 2] Get exam entity to be closed
	exam, err := uc.examRepository.GetExamByID(ctx, input.ExamID)
	if err != nil {
		return nil, err
	}

	if exam == nil {
		return nil, fmt.Errorf("exam with id %q does not exist", input.ExamID)
	}

	// [STEP 3] Verify that exam belongs to teacher
	if exam.ProfessorID != user.ID {
		return nil, fmt.Errorf("user does not have permissions to close this exam")
	}

	// [STEP 4] Read actual time (now)
	now := services.Now()

	// [STEP 5] Update exam entity with closing time
	exam, err = mapper.MapExamEndTimeInputToExamEntity(exam, now)
	if err != nil {
		return nil, err
	}

	// [STEP 6] Save updated exam entity
	exam, err = uc.examRepository.UpdateExam(ctx, exam)
	if err != nil {
		return nil, err
	}

	return exam, nil
}
