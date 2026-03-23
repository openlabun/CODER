package courses_usescases

import (
	"context"
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/course"

	userEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	repositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/course"
	userRepositoty "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
)

type GetCourseStudentsUseCase struct {
	courseRepository repositories.CourseRepository
	userRepository userRepositoty.UserRepository
}

func NewGetCourseStudentsUseCase(courseRepository repositories.CourseRepository, userRepository userRepositoty.UserRepository) *GetCourseStudentsUseCase {
	return &GetCourseStudentsUseCase{courseRepository: courseRepository, userRepository: userRepository}
}

func (uc *GetCourseStudentsUseCase) Execute(ctx context.Context, input dtos.GetCourseStudentsInput) ([]*userEntities.User, error) {
	// Check if course exists
	course, err := uc.courseRepository.GetCourseByID(ctx, input.CourseID)
	if err != nil {
		return nil, err
	}

	if course == nil {
		return nil, fmt.Errorf("course not found")
	}

	// Get students enrolled in course
	students, err := uc.courseRepository.GetStudentsByCourseID(ctx, input.CourseID)
	if err != nil {
		return nil, err
	}

	return students, nil
}