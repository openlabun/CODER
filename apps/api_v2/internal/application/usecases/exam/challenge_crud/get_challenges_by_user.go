package exam_usecases

import (
	"context"
	"fmt"

	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	
	user_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/user"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	repositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
)

type GetChallengesByUserUseCase struct {
	challengeRepository repositories.ChallengeRepository
	userRepository      userRepository.UserRepository
	examRepository      examRepository.ExamRepository
}

func NewGetChallengesByUserUseCase(challengeRepository repositories.ChallengeRepository, userRepository userRepository.UserRepository, examRepository examRepository.ExamRepository) *GetChallengesByUserUseCase {
	return &GetChallengesByUserUseCase{challengeRepository: challengeRepository, userRepository: userRepository, examRepository: examRepository}
}

func (uc *GetChallengesByUserUseCase) Execute(ctx context.Context, input dtos.GetChallengesByUserInput) ([]*Entities.Challenge, error) {
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

	if user.Role == user_constants.UserRoleStudent {
		return nil, fmt.Errorf("student users do not have access to challenges")
	}

	// [STEP 2] If exam id is provided, verify that exam exists and belongs to the user
	var exam *Entities.Exam
	if input.ExamID != nil {
		exam, err = uc.examRepository.GetExamByID(ctx, *input.ExamID)
		if err != nil {
			return nil, err
		}

		if exam == nil {
			return nil, fmt.Errorf("exam not found")
		}

		if exam.ProfessorID != user.ID {
			return nil, fmt.Errorf("user does not have access to the exam")
		}
	}

	// [STEP 3] Get challenges by user, if exam id is provided, filter by exam id
	challenges, err := uc.challengeRepository.GetChallengesByUserID(ctx, user.ID, input.ExamID)
	if err != nil {
		return nil, err
	}

	// [STEP 5] Return challenge details
	return challenges, nil
}