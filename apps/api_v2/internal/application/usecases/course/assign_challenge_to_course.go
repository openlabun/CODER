package courses_usescases

import (
	"context"
	"fmt"

	exam_repositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	user_repositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
)

type AssignChallengeToCourseUseCase struct {
	challengeRepository exam_repositories.ChallengeRepository
	userRepository      user_repositories.UserRepository
}

func NewAssignChallengeToCourseUseCase(challengeRepository exam_repositories.ChallengeRepository, userRepository user_repositories.UserRepository) *AssignChallengeToCourseUseCase {
	return &AssignChallengeToCourseUseCase{challengeRepository: challengeRepository, userRepository: userRepository}
}

func (uc *AssignChallengeToCourseUseCase) Execute(ctx context.Context, courseID, challengeID string) error {
	// [STEP 1] Verify user and get its role (only professor can assign)
	userEmail, err := services.UserEmailFromContext(ctx)
	if err != nil {
		return err
	}

	user, err := uc.userRepository.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return err
	}

	if user.Role != "professor" && user.Role != "admin" {
		return fmt.Errorf("only professors can assign challenges to courses")
	}

	// [STEP 2] Fetch challenge
	challenge, err := uc.challengeRepository.GetChallengeByID(ctx, challengeID)
	if err != nil {
		return err
	}
	if challenge == nil {
		return fmt.Errorf("challenge not found")
	}

	// [STEP 3] Update challenge with CourseID
	challenge.CourseID = courseID
	_, err = uc.challengeRepository.UpdateChallenge(ctx, challenge)
	return err
}
