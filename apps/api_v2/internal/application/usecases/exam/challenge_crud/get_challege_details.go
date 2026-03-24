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
	// [STEP 1] Verify user
	userEmail, err := services.UserEmailFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepository.GetUserByEmail(ctx, userEmail)
	if err != nil || user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// [STEP 2] Verify challenge exists
	challenge, err := uc.challengeRepository.GetChallengeByID(ctx, input.ChallengeID)
	if err != nil || challenge == nil {
		return nil, fmt.Errorf("challenge not found")
	}

	// [STEP 3] Access Control
	if user.Role == user_entities.UserRoleStudent {
		// Students can only see published challenges or those in an active exam
		if challenge.Status != Entities.ChallengeStatusPublished && challenge.ExamID == "" {
			return nil, fmt.Errorf("challenge is not accessible")
		}
	} else if user.Role == user_entities.UserRoleProfessor {
		// Professors should be able to see their own challenges
		// If it belongs to an exam, check exam ownership
		if challenge.ExamID != "" {
			exam, err := uc.examRepository.GetExamByID(ctx, challenge.ExamID)
			if err == nil && exam != nil && exam.ProfessorID != user.ID {
				return nil, fmt.Errorf("you do not have permission to view this exam's challenge")
			}
		}
	}

	return challenge, nil
}