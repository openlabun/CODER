package exam_usecases

import (
	"context"
	"fmt"
	
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	courseRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/course"
	user_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	courseEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
)

type GetExamDetailsUseCase struct {
	userRepository userRepository.UserRepository
	examRepository examRepository.ExamRepository
	courseRepository courseRepository.CourseRepository
}

func NewGetExamDetailsUseCase(userRepository userRepository.UserRepository, examRepository examRepository.ExamRepository, courseRepository courseRepository.CourseRepository) *GetExamDetailsUseCase {
	return &GetExamDetailsUseCase{
		userRepository: userRepository,
		examRepository: examRepository,
		courseRepository: courseRepository,
	}
}

func (uc *GetExamDetailsUseCase) Execute(ctx context.Context, input dtos.GetExamDetailsInput) (*Entities.Exam, error) {
	// [STEP 1] Verify user is teacher or student
	userEmail, err := services.UserEmailFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepository.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("user with email %q does not exist", userEmail)
	}

	role := user.Role

	// [STEP 2] Get exam entity with user provided exam ID
	exam, err := uc.examRepository.GetExamByID(ctx, input.ExamID)
	if err != nil {
		return nil, err
	}

	if exam == nil {
		return nil, fmt.Errorf("exam with id %q does not exist", input.ExamID)
	}

	// [STEP 3] If user is teacher and is not the owner or exam visibility is not "public" or "teachers"
	if role == user_entities.UserRoleProfessor && exam.ProfessorID != user.ID && exam.Visibility != Entities.VisibilityPublic && exam.Visibility != Entities.VisibilityTeachers {
		return nil, fmt.Errorf("user does not have permissions to view exam details")
	}

	// [STEP 4] If user is student and exam visibility is not "public"
	if role == user_entities.UserRoleStudent && exam.Visibility != Entities.VisibilityPublic {
		return nil, fmt.Errorf("user does not have permissions to view exam details")
	}

	// [STEP 5] If user is student, get its courses and verify that at least one of them is the course of the exam
	if role == user_entities.UserRoleStudent {
		if exam.Visibility != Entities.VisibilityCourse {
			return nil, fmt.Errorf("user does not have permissions to view exam details")
		}
		
		courses, err := uc.courseRepository.GetCoursesByStudentID(ctx, user.ID)
		if err != nil {
			return nil, err
		}

		if !studentInCourse(exam.CourseID, courses) {
			return nil, fmt.Errorf("user does not have permissions to view exam details")
		}
	}

	return exam, nil
}

func studentInCourse (courseID string, courses []*courseEntities.Course) bool {
    if courses == nil {
		return false
	}

	for _, course := range courses {
		if course == nil {
			continue
		}

		if course.ID == courseID {
			return true
		}
	}
	return false
}
