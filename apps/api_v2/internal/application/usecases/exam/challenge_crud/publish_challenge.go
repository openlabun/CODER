package exam_usecases

import (
	"context"
	"fmt"

	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	
	user_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	repositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	mapper "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam/mapper"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
)

type PublishChallengeUseCase struct {
	challengeRepository repositories.ChallengeRepository
	userRepository      userRepository.UserRepository
}

func NewPublishChallengeUseCase(challengeRepository repositories.ChallengeRepository, userRepository userRepository.UserRepository) *PublishChallengeUseCase {
	return &PublishChallengeUseCase{challengeRepository: challengeRepository, userRepository: userRepository}
}

func (uc *PublishChallengeUseCase) Execute(ctx context.Context, input dtos.PublishChallengeInput) (*Entities.Challenge, error) {
	// [STEP 1] Verify user is teacher and has permissions to publish an exam
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

	if user.Role != user_entities.UserRoleProfessor {
		return nil, fmt.Errorf("user does not have permissions to create an exam")
	}

	// [STEP 2] Verify challenge exists
	existingChallenge, err := uc.challengeRepository.GetChallengeByID(ctx, input.ChallengeID)
	if err != nil {
		return nil, err
	}

	if existingChallenge == nil {
		return nil, fmt.Errorf("challenge with id %q does not exist", input.ChallengeID)
	}

	// [STEP 3] Create challenge publish entity with user provided values
	challenge, err := mapper.MapPublishChallengeInputToChallengeEntity(existingChallenge)
	if err != nil {
		return nil, err
	}

	// [STEP 4] Save challenge with user provided values
	updatedChallenge, err := uc.challengeRepository.UpdateChallenge(ctx, challenge)
	if err != nil {
		return nil, err
	}
	
	return updatedChallenge, nil
}