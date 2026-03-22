package exam_usecases

import (
	"context"
	"fmt"
	
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	courseRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/course"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	user_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
)

type GetExamsByCourseUseCase struct {
	userRepository userRepository.UserRepository
	examRepository examRepository.ExamRepository
	courseRepository courseRepository.CourseRepository
}

func NewGetExamsByCourseUseCase(userRepository userRepository.UserRepository, examRepository examRepository.ExamRepository, courseRepository courseRepository.CourseRepository) *GetExamsByCourseUseCase {
	return &GetExamsByCourseUseCase{
		userRepository: userRepository,
		examRepository: examRepository,
		courseRepository: courseRepository,
	}
}

func (uc *GetExamsByCourseUseCase) Execute(ctx context.Context, input dtos.GetExamsByCourseInput) ([]*Entities.Exam, error) {
	// [STEP 1] Verify user is teacher
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

	// [STEP 2] Verify course exists and user is owner
	course, err := uc.courseRepository.GetCourseByID(ctx, input.CourseID)
	if err != nil {
		return nil, err
	}

	if course.ProfessorID != user.ID {
		return nil, fmt.Errorf("user is not the owner of the course")
	}

	// [STEP 3] Get exams for course
	exams, err := uc.examRepository.GetExamsByCourseID(ctx, input.CourseID)
	if err != nil {
		return nil, err
	}

	return exams, nil
}
