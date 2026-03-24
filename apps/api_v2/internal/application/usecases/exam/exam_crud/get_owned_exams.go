package exam_usecases

import (
	"context"
	"fmt"

	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	user_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	courseRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/course"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
)

type GetOwnedExamsUseCase struct {
	userRepository   userRepository.UserRepository
	examRepository   examRepository.ExamRepository
	courseRepository courseRepository.CourseRepository
}

func NewGetOwnedExamsUseCase(userRepository userRepository.UserRepository, examRepository examRepository.ExamRepository, courseRepository courseRepository.CourseRepository) *GetOwnedExamsUseCase {
	return &GetOwnedExamsUseCase{
		userRepository:   userRepository,
		examRepository:   examRepository,
		courseRepository: courseRepository,
	}
}

func (uc *GetOwnedExamsUseCase) Execute(ctx context.Context) ([]*Entities.Exam, error) {
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

	if user.Role != user_entities.UserRoleProfessor && user.Role != user_entities.UserRoleAdmin {
		return nil, fmt.Errorf("user does not have permissions to list owned exams")
	}

	if user.Role == user_entities.UserRoleProfessor {
		exams, err := uc.examRepository.GetExamsByTeacherID(ctx, user.ID)
		if err != nil {
			return nil, err
		}
		return exams, nil
	}

	courses, err := uc.courseRepository.GetAllCourses(ctx)
	if err != nil {
		return nil, err
	}

	seen := map[string]bool{}
	exams := make([]*Entities.Exam, 0)

	for _, course := range courses {
		if course == nil || course.ID == "" {
			continue
		}

		courseExams, err := uc.examRepository.GetExamsByCourseID(ctx, course.ID)
		if err != nil {
			return nil, err
		}

		for _, exam := range courseExams {
			if exam == nil || exam.ID == "" {
				continue
			}
			if seen[exam.ID] {
				continue
			}
			seen[exam.ID] = true
			exams = append(exams, exam)
		}
	}

	return exams, nil
}
