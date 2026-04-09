package exam_usecases

import (
	"context"
	"fmt"

	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	user_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"

	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
)

type GetPublicExamsUseCase struct {
	userRepository      userRepository.UserRepository
	examRepository      examRepository.ExamRepository
}

func NewGetPublicExamsUseCase(userRepository userRepository.UserRepository, examRepository examRepository.ExamRepository) *GetPublicExamsUseCase {
	return &GetPublicExamsUseCase{userRepository: userRepository, examRepository: examRepository}
}

func (uc *GetPublicExamsUseCase) Execute(ctx context.Context) ([]*Entities.Exam, error) {
	// [STEP 1] Verify user and get its role
	userEmail, err := services.UserEmailFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepository.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// [STEP 2] Get all published exams
	public_exams, err := uc.examRepository.GetPublicExams(ctx, string(Entities.VisibilityPublic))
	if err != nil {
		return nil, err
	}

	// [STEP 3] If user is a teacher, append exams with visibility for teachers
	if user.Role == user_entities.UserRoleProfessor {
		teacher_exams, err := uc.examRepository.GetPublicExams(ctx, string(Entities.VisibilityTeachers))
		if err != nil {
			return nil, err
		}
		public_exams = append(public_exams, teacher_exams...)
	}


	return public_exams, nil
}