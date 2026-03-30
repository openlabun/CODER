package exam_usecases

import (
	"context"
	"fmt"

	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	
	user_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	repositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
)

type ForkChallengeUseCase struct {
	challengeRepository repositories.ChallengeRepository
	userRepository      userRepository.UserRepository
	testCaseRepository repositories.TestCaseRepository
}

func NewForkChallengeUseCase(challengeRepository repositories.ChallengeRepository, userRepository userRepository.UserRepository, testCaseRepository repositories.TestCaseRepository) *ForkChallengeUseCase {
	return &ForkChallengeUseCase{challengeRepository: challengeRepository, userRepository: userRepository, testCaseRepository: testCaseRepository}
}

func (uc *ForkChallengeUseCase) Execute(ctx context.Context, input dtos.ForkChallengeInput) (*Entities.Challenge, error) {
	// [STEP 1] Verify user is teacher and has permissions to create an exam
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

	// [STEP 2] Get challenge to fork
	challengeToFork, err := uc.challengeRepository.GetChallengeByID(ctx, input.ChallengeID)
	if err != nil {
		return nil, err
	}
	if challengeToFork == nil {
		return nil, fmt.Errorf("challenge with id %q does not exist", input.ChallengeID)
	}

	// [STEP 3] Create challenge entity with user provided values
	forkedChallenge, err := services.ForkChallenge(ctx, *challengeToFork, user.ID, uc.challengeRepository, uc.testCaseRepository)
	if err != nil {
		return nil, err
	}
	if forkedChallenge == nil {
		return nil, fmt.Errorf("failed to fork challenge with id %q", input.ChallengeID)
	}

	// [STEP 4] Create challenge with user provided values
	createdChallenge, err := uc.challengeRepository.CreateChallenge(ctx, forkedChallenge)
	if err != nil {
		return nil, err
	}
	
	return createdChallenge, nil
}