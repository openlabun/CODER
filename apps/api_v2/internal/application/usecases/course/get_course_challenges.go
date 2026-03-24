package courses_usescases

import (
	"context"

	exam_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	user_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	exam_repositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	user_repositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
)

type GetCourseChallengesUseCase struct {
	challengeRepository exam_repositories.ChallengeRepository
	userRepository      user_repositories.UserRepository
}

func NewGetCourseChallengesUseCase(challengeRepository exam_repositories.ChallengeRepository, userRepository user_repositories.UserRepository) *GetCourseChallengesUseCase {
	return &GetCourseChallengesUseCase{challengeRepository: challengeRepository, userRepository: userRepository}
}

func (uc *GetCourseChallengesUseCase) Execute(ctx context.Context, courseID string) ([]*exam_entities.Challenge, error) {
	// [STEP 1] Verify user and get its role
	userEmail, err := services.UserEmailFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepository.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil // Or return unauthorized
	}

	role := user.Role

	// [STEP 2] Fetch challenges by CourseID
	challenges, err := uc.challengeRepository.GetChallengesByCourseID(ctx, courseID)
	if err != nil {
		return nil, err
	}

	// [STEP 3] Filter to only return published ones for students
	if role == user_entities.UserRoleStudent {
		challenges = filterPublishedChallenges(challenges)
	}

	return challenges, nil
}

func filterPublishedChallenges(challenges []*exam_entities.Challenge) []*exam_entities.Challenge {
	visibleChallenges := []*exam_entities.Challenge{}
	for _, challenge := range challenges {
		if challenge.Status == exam_entities.ChallengeStatusPublished {
			visibleChallenges = append(visibleChallenges, challenge)
		}
	}
	return visibleChallenges
}
