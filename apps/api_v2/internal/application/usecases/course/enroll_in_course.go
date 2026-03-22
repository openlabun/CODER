package courses_usescases

import (
	"context"
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/course"
	mapper "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/course/mapper"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"
	repositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/course"
	userRepositoty "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
)

type EnrollInCourseUseCase struct {
	courseRepository repositories.CourseRepository
	userRepository userRepositoty.UserRepository
}

func NewEnrollInCourseUseCase(courseRepository repositories.CourseRepository, userRepository userRepositoty.UserRepository) *EnrollInCourseUseCase {
	return &EnrollInCourseUseCase{courseRepository: courseRepository, userRepository: userRepository}
}

func (uc *EnrollInCourseUseCase) Execute(ctx context.Context, input dtos.EnrolledInCourseInput) (*Entities.Course, error) {
	// Check if student exists
	student, err := uc.userRepository.GetUserByID(ctx, input.StudentID)
	if err != nil {
		return nil, err
	}

	if student == nil {
		return nil, fmt.Errorf("student not found")
	}

	// Check if course exists
	course, err := uc.courseRepository.GetCourseByID(ctx, input.CourseID)
	if err != nil {
		return nil, err
	}

	if course == nil {
		return nil, fmt.Errorf("course not found")
	}

	// Check if student is already enrolled in course
	enrollment, err := uc.courseRepository.GetCoursesByStudentID(ctx, input.StudentID)
	if err != nil {
		return nil, err
	}

	if studentAlreadyEnrolled(enrollment, input.CourseID) {
		return nil, fmt.Errorf("student is already enrolled in this course")
	}

	// Enroll student in course
	mapper.MapCourseStudentInputToEntity(input)
	err = uc.courseRepository.AddStudentToCourse(ctx, input.CourseID, input.StudentID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func studentAlreadyEnrolled(enrollments []*Entities.Course, courseID string) bool {
	for _, enrollment := range enrollments {
		if enrollment.ID == courseID {
			return true
		}
	}

	return false
}