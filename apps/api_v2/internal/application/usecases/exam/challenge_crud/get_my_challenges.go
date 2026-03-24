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

type GetMyChallengesUseCase struct {
	challengeRepository examRepository.ChallengeRepository
	userRepository      userRepository.UserRepository
	courseRepository    courseRepository.CourseRepository
}

func NewGetMyChallengesUseCase(challengeRepository examRepository.ChallengeRepository, userRepository userRepository.UserRepository, courseRepository courseRepository.CourseRepository) *GetMyChallengesUseCase {
	return &GetMyChallengesUseCase{
		challengeRepository: challengeRepository,
		userRepository:      userRepository,
		courseRepository:    courseRepository,
	}
}

func (uc *GetMyChallengesUseCase) Execute(ctx context.Context) ([]*Entities.Challenge, error) {
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

	var courseIDs []string
	switch user.Role {
	case user_entities.UserRoleStudent:
		courses, err := uc.courseRepository.GetCoursesByStudentID(ctx, user.ID)
		if err != nil {
			return nil, err
		}
		courseIDs = make([]string, 0, len(courses))
		for _, course := range courses {
			if course != nil && course.ID != "" {
				courseIDs = append(courseIDs, course.ID)
			}
		}
	case user_entities.UserRoleProfessor:
		courses, err := uc.courseRepository.GetCoursesByTeacherID(ctx, user.ID)
		if err != nil {
			return nil, err
		}
		courseIDs = make([]string, 0, len(courses))
		for _, course := range courses {
			if course != nil && course.ID != "" {
				courseIDs = append(courseIDs, course.ID)
			}
		}
	case user_entities.UserRoleAdmin:
		courses, err := uc.courseRepository.GetAllCourses(ctx)
		if err != nil {
			return nil, err
		}
		courseIDs = make([]string, 0, len(courses))
		for _, course := range courses {
			if course != nil && course.ID != "" {
				courseIDs = append(courseIDs, course.ID)
			}
		}
	default:
		return nil, fmt.Errorf("unsupported user role")
	}

	if len(courseIDs) == 0 {
		return []*Entities.Challenge{}, nil
	}

	seen := map[string]bool{}
	challenges := make([]*Entities.Challenge, 0)

	for _, courseID := range courseIDs {
		items, err := uc.challengeRepository.GetChallengesByCourseID(ctx, courseID)
		if err != nil {
			return nil, err
		}

		for _, challenge := range items {
			if challenge == nil || challenge.ID == "" {
				continue
			}
			if seen[challenge.ID] {
				continue
			}
			seen[challenge.ID] = true
			challenges = append(challenges, challenge)
		}
	}

	if user.Role == user_entities.UserRoleStudent {
		challenges = filterPublishedChallenges(challenges)
	}

	return challenges, nil
}
