package courses_usescases

import (
	"context"
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/course"
	mapper "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/course/mapper"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"
	user_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	repositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/course"
	repositories_user "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
)

type UpdateCourseUseCase struct {
	userRepository   repositories_user.UserRepository
	courseRepository repositories.CourseRepository
}

func NewUpdateCourseUseCase(courseRepository repositories.CourseRepository, userRepository repositories_user.UserRepository) *UpdateCourseUseCase {
	return &UpdateCourseUseCase{courseRepository: courseRepository, userRepository: userRepository}
}

func (uc *UpdateCourseUseCase) Execute(ctx context.Context, input dtos.UpdateCourseInput) (*Entities.Course, error) {
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

	// Get original course entity with user provided values
	originalCourse, err := uc.courseRepository.GetCourseByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if originalCourse == nil {
		return nil, fmt.Errorf("course not found")
	}

	// Create course entity with user provided values
	course, err := mapper.MapUpdateCourseInputToCourseEntity(originalCourse, input)
	if err != nil {
		return nil, err
	}

	// Update course with user provided values
	updatedCourse, err := uc.courseRepository.UpdateCourse(ctx, course)
	if err != nil {
		return nil, err
	}

	return updatedCourse, nil
}
