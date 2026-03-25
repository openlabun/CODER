package courses_usescases

import (
	"context"
	"fmt"

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

func (uc *GetOwnedCoursesUseCase) Execute(ctx context.Context) ([]*Entities.Course, error) {
	// Verify user is a teacher
	userEmail, err := services.UserEmailFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepository.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return nil, err
	}

	if user.Role != user_entities.UserRoleProfessor {
		return nil, fmt.Errorf("user does not have permissions to view owned courses")
	}

	// Get owned courses for the teacher
	courses, err := uc.courseRepository.GetCoursesByTeacherID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	if courses == nil {
		return []*Entities.Course{}, nil
	}

	return courses, nil
}
