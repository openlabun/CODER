package exam_usecases

import (
	"context"
	"fmt"

	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	
	user_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	repositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
)

type GetChallengeDetailsUseCase struct {
	challengeRepository repositories.ChallengeRepository
	userRepository      userRepository.UserRepository
	examRepository      examRepository.ExamRepository
}

func NewGetChallengeDetailsUseCase(challengeRepository repositories.ChallengeRepository, userRepository userRepository.UserRepository, examRepository examRepository.ExamRepository) *GetChallengeDetailsUseCase {
	return &GetChallengeDetailsUseCase{challengeRepository: challengeRepository, userRepository: userRepository, examRepository: examRepository}
}

func (uc *GetChallengeDetailsUseCase) Execute(ctx context.Context, input dtos.GetChallengeDetailsInput) (*Entities.Challenge, error) {
	// [STEP 1] Verify user is teacher
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

	// [STEP 2] Verify challenge exists
	challenge, err := uc.challengeRepository.GetChallengeByID(ctx, input.ChallengeID)
	if err != nil {
		return nil, err
	}

	if challenge == nil {
		return nil, fmt.Errorf("challenge with id %s does not exist", input.ChallengeID)
	}

	// [STEP 3] Verify that challenge belongs to teacher
	if challenge.UserID != user.ID {
		return nil, fmt.Errorf("user does not have permissions to access this challenge")
	}

	return challenge, nil
}