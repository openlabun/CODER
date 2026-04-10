package courses_usescases

import (
	"context"
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/course"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"

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
	// [STEP 1] Check if student exists using ID or email, fallback to authenticated user
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
		// Fallback to authenticated user
		userEmail, err := services.UserEmailFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("student_id or student_email must be provided, or user must be authenticated")
		}
		student_, err := uc.userRepository.GetUserByEmail(ctx, userEmail)
		if err != nil {
			return nil, err
		}
		student = student_
	}

	if student == nil {
		return nil, fmt.Errorf("student not found")
	}

	// [STEP 2] Check if course exists
	var course *Entities.Course
	var err error

	if input.CourseID != "" {
		course, err = uc.courseRepository.GetCourseByID(ctx, input.CourseID)
	} else if input.EnrollmentCode != "" {
		course, err = uc.courseRepository.GetCourseByEnrollmentCode(ctx, input.EnrollmentCode)
	} else {
		return nil, fmt.Errorf("course_id or enrollment_code must be provided")
	}

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

	if studentAlreadyEnrolled(enrollment, course.ID) {
		return nil, fmt.Errorf("student is already enrolled in this course")
	}

	// [STEP 4] Enroll student in course
	err = uc.courseRepository.AddStudentToCourse(ctx, course.ID, student.ID)
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