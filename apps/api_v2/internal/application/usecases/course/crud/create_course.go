package courses_usescases

import (
	"context"
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/course"
	mapper "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/course/mapper"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"
	user_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/user"
	repositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/course"
	repositories_user "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
)

type CreateCourseUseCase struct {
	userRepository   repositories_user.UserRepository
	courseRepository repositories.CourseRepository
}

func NewCreateCourseUseCase(courseRepository repositories.CourseRepository, userRepository repositories_user.UserRepository) *CreateCourseUseCase {
	return &CreateCourseUseCase{courseRepository: courseRepository, userRepository: userRepository}
}

func (uc *CreateCourseUseCase) Execute(ctx context.Context, input dtos.CreateCourseInput) (*Entities.Course, error) {
	// Verify user is teacher and has permissions to create a course
	userEmail, err := services.UserEmailFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepository.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return nil, err
	}

	if user.Role != user_constants.UserRoleProfessor {
		return nil, fmt.Errorf("user does not have permissions to create a course")
	}

	// Create course entity with user provided values
	course, err := mapper.MapCreateCourseInputToCourseEntity(user.ID, input)
	if err != nil {
		return nil, err
	}

	// Create course with user provided values
	createdCourse, err := uc.courseRepository.CreateCourse(ctx, course)
	if err != nil {
		return nil, err
	}

	return createdCourse, nil
}
