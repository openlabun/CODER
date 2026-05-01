package submission_usecases

import (
	"context"
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"
	"github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission/mapper"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	user_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/user"
	examScoreRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
)

type GetExamScoringUseCase struct {
	userRepository userRepository.UserRepository
	examScoreRepository examScoreRepository.ExamScoreRepository
}

func NewGetExamScoringUseCase(userRepository userRepository.UserRepository, examScoreRepository examScoreRepository.ExamScoreRepository) *GetExamScoringUseCase {
	return &GetExamScoringUseCase{
		userRepository: userRepository,
		examScoreRepository: examScoreRepository,
	}
}

func (uc *GetExamScoringUseCase) Execute(ctx context.Context, input dtos.GetExamScoringInput) (*dtos.ExamScoringDTO, error) {
	// [STEP 1] Verify user is student and has permissions to submit
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

	// [STEP 2] Verify user is teacher or student who took the exam
	if user.Role != user_constants.UserRoleProfessor && user.ID != input.UserID {
		return nil, fmt.Errorf("User %q does not have permissions to view this exam scoring", user.Username)
	}

	// [STEP 3] Get exam scores by examID
	examScores, err := uc.examScoreRepository.GetExamScores(ctx, &input.ExamID, &input.UserID)
	if err != nil {
		return nil, err
	}

	if len(examScores) == 0 {
		return nil, fmt.Errorf("no exam scores found for examID %q and userID %q", input.ExamID, input.UserID)
	}

	// [STEP 4] Get best score from all tries
	best_score := getBestScore(examScores)

	// [STEP 5] Map exam scores into output DTO
	examScoringDTO, err := mapper.MapExamScoreEntitiesToExamScoringDTO(examScores, best_score)
	if err != nil {
		return nil, err
	}

	return examScoringDTO, nil
}

func getBestScore(examScores []*Entities.ExamScore) int {
	if len(examScores) == 0 {
		return 0
	}
	bestScore := examScores[0].Score
	for _, score := range examScores {
		if score.Score > bestScore {
			bestScore = score.Score
		}
	}
	return bestScore
}