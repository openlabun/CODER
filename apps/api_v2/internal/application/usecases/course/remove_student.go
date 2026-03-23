package courses_usescases

import (
	"context"
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/course"

	repositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/course"
	userRepositoty "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
)

type RemoveStudentFromCourseUseCase struct {
	courseRepository repositories.CourseRepository
	userRepository userRepositoty.UserRepository
}

func NewRemoveStudentFromCourseUseCase(courseRepository repositories.CourseRepository, userRepository userRepositoty.UserRepository) *RemoveStudentFromCourseUseCase {
	return &RemoveStudentFromCourseUseCase{courseRepository: courseRepository, userRepository: userRepository}
}

func (uc *RemoveStudentFromCourseUseCase) Execute(ctx context.Context, input dtos.RemoveStudentFromCourseInput) error {
	// Check if student exists
	student, err := uc.userRepository.GetUserByID(ctx, input.StudentID)
	if err != nil {
		return err
	}

	if student == nil {
		return fmt.Errorf("student not found")
	}

	// Check if course exists
	course, err := uc.courseRepository.GetCourseByID(ctx, input.CourseID)
	if err != nil {
		return err
	}

	if course == nil {
		return fmt.Errorf("course not found")
	}

	// Check if student is already enrolled in course
	enrollment, err := uc.courseRepository.GetCoursesByStudentID(ctx, input.StudentID)
	if err != nil {
		return err
	}

	if !studentAlreadyEnrolled(enrollment, input.CourseID) {
		return fmt.Errorf("student is not enrolled in this course")
	}

	// Remove student from course
	err = uc.courseRepository.RemoveStudentFromCourse(ctx, input.CourseID, input.StudentID)
	if err != nil {
		return err
	}

	return nil
}