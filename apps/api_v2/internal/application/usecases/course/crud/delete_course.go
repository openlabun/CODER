package courses_usescases

import (
	"context"
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/course"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"

	user_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	repositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/course"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
)

type DeleteCourseUseCase struct {
	courseRepository repositories.CourseRepository
	userRepository   userRepository.UserRepository
}

func NewDeleteCourseUseCase(courseRepository repositories.CourseRepository, userRepository userRepository.UserRepository) *DeleteCourseUseCase {
	return &DeleteCourseUseCase{courseRepository: courseRepository, userRepository: userRepository}
}

func (uc *DeleteCourseUseCase) Execute(ctx context.Context, input dtos.DeleteCourseInput) error {
	// Verify user is teacher and has permissions to create a course
	userEmail, err := services.UserEmailFromContext(ctx)
	if err != nil {
		return err
	}

	user, err := uc.userRepository.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return err
	}

	if user.Role != user_entities.UserRoleProfessor {
		return fmt.Errorf("user does not have permissions to create a course")
	}

	// Delete course with user provided values
	err = uc.courseRepository.DeleteCourse(ctx, input.CourseID)
	if err != nil {
		return err
	}

	return nil
}
