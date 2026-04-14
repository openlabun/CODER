package exam_usecases

import (
	"context"
	"fmt"

	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	user_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/user"
	repositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	domain_services "github.com/openlabun/CODER/apps/api_v2/internal/domain/services"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	mapper "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam/mapper"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
)

type CreateChallengeUseCase struct {
	challengeRepository  repositories.ChallengeRepository
	ioVariableRepository repositories.IOVariableRepository
	userRepository       userRepository.UserRepository
}

func NewCreateChallengeUseCase(challengeRepository repositories.ChallengeRepository, ioVariableRepository repositories.IOVariableRepository, userRepository userRepository.UserRepository) *CreateChallengeUseCase {
	return &CreateChallengeUseCase{challengeRepository: challengeRepository, ioVariableRepository: ioVariableRepository, userRepository: userRepository}
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

	if user == nil {
		return nil, fmt.Errorf("user with email %q does not exist", userEmail)
	}

	if user.Role != user_constants.UserRoleProfessor {
		return nil, fmt.Errorf("user does not have permissions to create an exam")
	}

	// [STEP 2] Set userID in input to create challenge with the authenticated user as owner
	input.UserID = user.ID

	// [STEP 3] Create challenge entity with user provided values
	challenge, err := mapper.MapCreateChallengeInputToChallengeEntity(input)
	if err != nil {
		return nil, err
	}

	// [STEP 4] Create challenge with user provided values
	createdChallenge, err := domain_services.CreateChallenge(ctx, challenge, uc.challengeRepository, uc.ioVariableRepository)
	if err != nil {
		return nil, err
	}

	return createdChallenge, nil
}
