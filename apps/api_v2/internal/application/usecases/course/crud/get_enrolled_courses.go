package courses_usescases

import (
	"context"
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/course"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"
	user_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	repositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/course"
	userRepositoty "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
)

type GetEnrolledCoursesUseCase struct {
	courseRepository repositories.CourseRepository
	userRepository   userRepositoty.UserRepository
}

func NewGetEnrolledCoursesUseCase(courseRepository repositories.CourseRepository, userRepository userRepositoty.UserRepository) *GetEnrolledCoursesUseCase {
	return &GetEnrolledCoursesUseCase{courseRepository: courseRepository, userRepository: userRepository}
}

func (uc *GetEnrolledCoursesUseCase) Execute(ctx context.Context, input dtos.GetEnrolledCoursesInput) ([]*Entities.Course, error) {
	// Verify user is a student
	userEmail, err := services.UserEmailFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepository.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return nil, err
	}

	if user.Role != user_entities.UserRoleStudent {
		return nil, fmt.Errorf("user does not have permissions to view enrolled courses")
	}

	// Get enrolled courses for the student
	courses, err := uc.courseRepository.GetCoursesByStudentID(ctx, input.StudentID)
	if err != nil {
		return nil, err
	}

	if courses == nil {
		return nil, fmt.Errorf("no enrolled courses found")
	}

	return courses, nil
}
