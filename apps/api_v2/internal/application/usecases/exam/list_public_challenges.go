package exam_usecases

import (
	"context"
	"fmt"

	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	
	user_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
)

type GetPublicChallengesUseCase struct {
	challengeRepository examRepository.ChallengeRepository
	userRepository      userRepository.UserRepository
}

func NewGetPublicChallengesUseCase(challengeRepository examRepository.ChallengeRepository, userRepository userRepository.UserRepository) *GetPublicChallengesUseCase {
	return &GetPublicChallengesUseCase{challengeRepository: challengeRepository, userRepository: userRepository}
}

func (uc *GetPublicChallengesUseCase) Execute(ctx context.Context, input dtos.GetPublicChallengesInput) ([]*Entities.Challenge, error) {
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
		return nil, fmt.Errorf("user not found")
	}

	if user.Role == user_entities.UserRoleStudent {
		return nil, fmt.Errorf("students are not allowed to access challenges repositories")
	}

	// [STEP 2] Get all published exams
	status := string(Entities.ChallengeStatusPublished)
	public_exams, err := uc.challengeRepository.GetChallenges(ctx, &status, input.Tag, input.Difficulty)
	if err != nil {
		return nil, err
	}

	return public_exams, nil
}