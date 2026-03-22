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

type GetCourseDetailsUseCase struct {
	courseRepository repositories.CourseRepository
	userRepository   userRepositoty.UserRepository
}

func NewGetCourseDetailsUseCase(courseRepository repositories.CourseRepository, userRepository userRepositoty.UserRepository) *GetCourseDetailsUseCase {
	return &GetCourseDetailsUseCase{courseRepository: courseRepository, userRepository: userRepository}
}

func (uc *GetCourseDetailsUseCase) Execute(ctx context.Context, input dtos.GetCourseDetailsInput) (*Entities.Course, error) {
	// Verify user is teacher and has permissions to create a course
	userEmail, err := services.UserEmailFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepository.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return nil, err
	}

	if user.Role != user_entities.UserRoleProfessor {
		return nil, fmt.Errorf("user does not have permissions to create a course")
	}

	// Get course details with user provided values
	course, err := uc.courseRepository.GetCourseByID(ctx, input.CourseID)
	if err != nil {
		return nil, err
	}

	if course == nil {
		return nil, fmt.Errorf("course not found")
	}

	return course, nil
}
