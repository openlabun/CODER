package courses_usescases

import (
	"context"
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/course"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"
	userEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
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
	// [STEP 1] Check if student exists using ID or email
	var student *userEntities.User
	if input.StudentID != nil {
		student_, err := uc.userRepository.GetUserByID(ctx, *input.StudentID)
		if err != nil {
			return nil, err
		}
		student = student_
	} else if input.StudentEmail != nil {
		student_, err := uc.userRepository.GetUserByEmail(ctx, *input.StudentEmail)
		if err != nil {
			return nil, err
		}
		student = student_
	} else {
		return nil, fmt.Errorf("student_id or student_email must be provided")
	}

	if student == nil {
		return nil, fmt.Errorf("student not found")
	}

	// [STEP 2] Check if course exists
	course, err := uc.courseRepository.GetCourseByID(ctx, input.CourseID)
	if err != nil {
		return nil, err
	}

	if course == nil {
		return nil, fmt.Errorf("course not found")
	}

	// [STEP 3] Check if student is already enrolled in course
	enrollment, err := uc.courseRepository.GetCoursesByStudentID(ctx, student.ID)
	if err != nil {
		return nil, err
	}

	if studentAlreadyEnrolled(enrollment, input.CourseID) {
		return nil, fmt.Errorf("student is already enrolled in this course")
	}

	// [STEP 4] Enroll student in course
	err = uc.courseRepository.AddStudentToCourse(ctx, input.CourseID, student.ID)
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