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
	mapper "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam/mapper"
)

type CreateExamUseCase struct {
	userRepository userRepository.UserRepository
	courseRepository courseRepository.CourseRepository
	examRepository examRepository.ExamRepository
}

func NewCreateExamUseCase(userRepository userRepository.UserRepository, courseRepository courseRepository.CourseRepository, examRepository examRepository.ExamRepository) *CreateExamUseCase {
	return &CreateExamUseCase{
		userRepository: userRepository,
		courseRepository: courseRepository,
		examRepository: examRepository,
	}
}

func (uc *CreateExamUseCase) Execute(ctx context.Context, input dtos.CreateExamInput) (*Entities.Exam, error) {
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

	if user.Role != user_entities.UserRoleProfessor {
		return nil, fmt.Errorf("user does not have permissions to create an exam")
	}

	// [STEP 2] Validate course association if provided
	if input.CourseID != nil {
		course, err := uc.courseRepository.GetCourseByID(ctx, *input.CourseID)
		if err != nil {
			return nil, err
		}
		
		if course == nil {
			return nil, fmt.Errorf("course with ID %q does not exist", *input.CourseID)
		}
	}

	// [STEP 2] Create exam entity with user provided values
	exam, err := mapper.MapCreateExamInputToExamEntity(input)
	if err != nil {
		return nil, err
	}


	// [STEP 3] Create exam with user provided values
	exam, err = uc.examRepository.CreateExam(ctx, exam)
	if err != nil {
		return nil, err
	}

	return exam, nil
}
