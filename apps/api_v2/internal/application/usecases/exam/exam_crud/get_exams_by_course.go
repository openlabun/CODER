package exam_usecases

import (
	"context"
	"fmt"

	constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/exam"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	user_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/user"
	courseRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/course"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"

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
	// [STEP 1] Verify user
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

	// [STEP 2] Verify course exists
	course, err := uc.courseRepository.GetCourseByID(ctx, input.CourseID)
	if err != nil {
		return nil, err
	}

	if course == nil {
		return nil, fmt.Errorf("course with id %q does not exist", input.CourseID)
	}

	// [STEP 3] Get exams for course
	exams, err := uc.examRepository.GetExamsByCourseID(ctx, input.CourseID)
	if err != nil {
		return nil, err
	}

	// [STEP 4] If user is student, filter exams by visibility
	if user.Role == user_constants.UserRoleStudent {
		filteredExams := []*Entities.Exam{}
		for _, exam := range exams {
			if exam.Visibility == constants.VisibilityCourse {
				filteredExams = append(filteredExams, exam)
			}
		}
		
		exams = filteredExams
	}

	return exams, nil
}
