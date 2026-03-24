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

type GetOwnedCoursesUseCase struct {
	courseRepository repositories.CourseRepository
	userRepository   userRepositoty.UserRepository
}

func NewGetOwnedCoursesUseCase(courseRepository repositories.CourseRepository, userRepository userRepositoty.UserRepository) *GetOwnedCoursesUseCase {
	return &GetOwnedCoursesUseCase{courseRepository: courseRepository, userRepository: userRepository}
}

func (uc *GetOwnedCoursesUseCase) Execute(ctx context.Context, input dtos.GetOwnedCoursesInput) ([]*Entities.Course, error) {
	// Verify user is a teacher
	userEmail, err := services.UserEmailFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepository.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return nil, err
	}

	// Allow both professor and teacher roles for compatibility
	if user.Role != user_entities.UserRoleProfessor && user.Role != "teacher" {
		return nil, fmt.Errorf("user role '%s' does not have permissions to view owned courses", user.Role)
	}

	// Use user ID from context if TeacherID is not provided in query
	teacherID := input.TeacherID
	if teacherID == "" {
		teacherID = user.ID
	}

	// Get owned courses for the teacher
	courses, err := uc.courseRepository.GetCoursesByTeacherID(ctx, teacherID)
	if err != nil {
		return nil, err
	}

	// Return empty list instead of nil if no courses found
	if courses == nil {
		return []*Entities.Course{}, nil
	}

	return courses, nil
}
