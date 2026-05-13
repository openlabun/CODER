package submission_usecases

import (
	"fmt"
	"context"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
	mapper "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission/mapper"

	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	user_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/user"
	examScoreRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
)

type GetUsersScoresByExamUseCase struct {
	userRepository userRepository.UserRepository
	examScoreRepository examScoreRepository.ExamScoreRepository
	examRepository examRepository.ExamRepository
}

func NewGetUsersScoresByExamUseCase(userRepository userRepository.UserRepository, examScoreRepository examScoreRepository.ExamScoreRepository, examRepository examRepository.ExamRepository) *GetUsersScoresByExamUseCase {
	return &GetUsersScoresByExamUseCase{
		userRepository: userRepository,
		examScoreRepository: examScoreRepository,
		examRepository: examRepository,
	}
}

func (uc *GetUsersScoresByExamUseCase) Execute(ctx context.Context, input dtos.GetUserScoresByExamInput) (*dtos.UsersExamScoresOutputDTO, error) {
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

	// [STEP 2] Verify user is teacher
	if user.Role != user_constants.UserRoleProfessor {
		return nil, fmt.Errorf("User %q does not have permissions to view this exam scoring", user.Username)
	}

	// [STEP 3] Get exam and validate it exists and its owned by teacher
	exam, err := uc.examRepository.GetExamByID(ctx, input.ExamID)
	if err != nil {
		return nil, err
	}

	if exam == nil {
		return nil, fmt.Errorf("Exam with ID %q does not exist", input.ExamID)
	}

	if exam.ProfessorID != user.ID {
		return nil, fmt.Errorf("User %q does not have permissions to view this exam scoring", user.Username)
	}

	// [STEP 4] Get exams scores by userID
	examScores, err := uc.examScoreRepository.GetExamScores(ctx, &input.ExamID, nil)
	if err != nil {
		return nil, err
	}
	if len(examScores) == 0 {
		return nil, fmt.Errorf("no exam scores found for examID %q", input.ExamID)
	}

	// [STEP 5] Aggregate scores by user
	aggregatedScores := services.AggregateScoresByUser(examScores)

	// [STEP 6] Map to ExamScoringDTO and UserExamScoresOutputDTO
	var dtos []*dtos.UserExamScoresOutputDTO
	for userID, scores := range aggregatedScores {
		bestScore := services.CalculateBestScore(scores)
		examScoringDTO, err := mapper.MapExamScoreEntitiesToExamScoringDTO(scores, bestScore)
		if err != nil {
			return nil, err
		}

		userExamScoresDTO, err := mapper.MapExamScoresDTOtoUserExamScoresOutputDTO(examScoringDTO, userID)
		if err != nil {
			return nil, err
		}

		dtos = append(dtos, userExamScoresDTO)
	}

	// [STEP 7] Map to output DTO and return
	usersExamScoresDTO, err := mapper.MapUserExamScoresOutputDTOToUsersExamScoresOutputDTO(dtos, input.ExamID)
	if err != nil {
		return nil, err
	}

	return usersExamScoresDTO, nil
}