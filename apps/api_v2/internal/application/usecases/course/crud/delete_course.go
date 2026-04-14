package courses_usescases

import (
	"context"
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/course"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"

	domain_services "github.com/openlabun/CODER/apps/api_v2/internal/domain/services"
	user_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/user"
	repositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/course"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
)

type DeleteCourseUseCase struct {
	courseRepository repositories.CourseRepository
	userRepository   userRepository.UserRepository
	examRepository   examRepository.ExamRepository
	examScoreRepository examRepository.ExamScoreRepository
	examItemRepository examRepository.ExamItemRepository
	examItemScoreRepository examRepository.ExamItemScoreRepository
}

func NewDeleteCourseUseCase(courseRepository repositories.CourseRepository, userRepository userRepository.UserRepository, examRepository examRepository.ExamRepository, examItemRepository examRepository.ExamItemRepository, examScoreRepository examRepository.ExamScoreRepository, examItemScoreRepository examRepository.ExamItemScoreRepository) *DeleteCourseUseCase {
	return &DeleteCourseUseCase{courseRepository: courseRepository, userRepository: userRepository, examRepository: examRepository, examScoreRepository: examScoreRepository, examItemRepository: examItemRepository, examItemScoreRepository: examItemScoreRepository}
}

func (uc *DeleteCourseUseCase) Execute(ctx context.Context, input dtos.DeleteCourseInput) error {
	// [STEP 1] Verify user is teacher and has permissions to create a course
	userEmail, err := services.UserEmailFromContext(ctx)
	if err != nil {
		return err
	}

	user, err := uc.userRepository.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return err
	}

	if user.Role != user_constants.UserRoleProfessor {
		return fmt.Errorf("user does not have permissions to create a course")
	}

	// [STEP 2] Delete course with user provided values
	err = domain_services.RemoveCourse(ctx, input.CourseID, uc.courseRepository, uc.examRepository, uc.examItemRepository, uc.examScoreRepository, uc.examItemScoreRepository)
	if err != nil {
		return err
	}

	return nil
}
