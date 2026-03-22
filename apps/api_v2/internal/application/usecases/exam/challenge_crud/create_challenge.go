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

type CreateChallengeUseCase struct {
	challengeRepository repositories.ChallengeRepository
	userRepository      userRepository.UserRepository
}

func NewCreateChallengeUseCase(challengeRepository repositories.ChallengeRepository, userRepository userRepository.UserRepository) *CreateChallengeUseCase {
	return &CreateChallengeUseCase{challengeRepository: challengeRepository, userRepository: userRepository}
}

func (uc *CreateChallengeUseCase) Execute(ctx context.Context, input dtos.CreateChallengeInput) (*Entities.Challenge, error) {
	// [STEP 1] Verify user is teacher and has permissions to create an exam
	userEmail, err := services.UserEmailFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepository.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return nil, err
	}

	if user.Role != user_entities.UserRoleProfessor {
		return nil, fmt.Errorf("user does not have permissions to create an exam")
	}

	// [STEP 2] Create challenge entity with user provided values
	challenge, err := mapper.MapCreateChallengeInputToChallengeEntity(input)
	if err != nil {
		return nil, err
	}

	// [STEP 3] Create challenge with user provided values
	createdChallenge, err := uc.challengeRepository.CreateChallenge(ctx, challenge)
	if err != nil {
		return nil, err
	}
	
	return createdChallenge, nil
}