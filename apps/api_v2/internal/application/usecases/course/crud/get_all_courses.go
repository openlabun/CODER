package courses_usescases

import (
	"context"
	"fmt"
	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/course"

	repositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/course"
	userRepositoty "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
)

type GetAllCoursesUseCase struct {
	courseRepository repositories.CourseRepository
	userRepository   userRepositoty.UserRepository
}

func NewGetAllCoursesUseCase(courseRepository repositories.CourseRepository, userRepository userRepositoty.UserRepository) *GetAllCoursesUseCase {
	return &GetAllCoursesUseCase{courseRepository: courseRepository, userRepository: userRepository}
}

func (uc *GetAllCoursesUseCase) Execute(ctx context.Context) ([]*dtos.CourseBrowseItem, error) {
	courses, err := uc.courseRepository.GetAllCourses(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*dtos.CourseBrowseItem, 0, len(courses))
	for _, c := range courses {
		professorName := "Unknown"
		if c.ProfessorID != "" {
			prof, err := uc.userRepository.GetUserByID(ctx, c.ProfessorID)
			if err == nil && prof != nil {
				professorName = prof.Username
			}
		}

		formattedPeriod := "Unknown"
		if c.Period != nil {
			formattedPeriod = fmt.Sprintf("%d-%s", c.Period.Year, c.Period.Semester)
		}

		result = append(result, &dtos.CourseBrowseItem{
			ID:            c.ID,
			Name:          c.Name,
			Code:          c.Code,
			Period:        formattedPeriod,
			ProfessorID:   c.ProfessorID,
			ProfessorName: professorName,
			CreatedAt:     c.CreatedAt.String(),
		})
	}

	return result, nil
}
