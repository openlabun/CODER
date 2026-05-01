package submission_usecases

import (
	"context"
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"
	"github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission/mapper"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"

	user_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/user"
	examEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	examScoreRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
)

type GetExamsScoresByUserUseCase struct {
	userRepository userRepository.UserRepository
	examScoreRepository examScoreRepository.ExamScoreRepository
	examRepository examRepository.ExamRepository
}

func NewGetExamsScoresByUserUseCase(userRepository userRepository.UserRepository, examScoreRepository examScoreRepository.ExamScoreRepository, examRepository examRepository.ExamRepository) *GetExamsScoresByUserUseCase {
	return &GetExamsScoresByUserUseCase{
		userRepository: userRepository,
		examScoreRepository: examScoreRepository,
		examRepository: examRepository,
	}
}

func (uc *GetExamsScoresByUserUseCase) Execute(ctx context.Context, input dtos.GetExamScoresByUserInput) (*dtos.UserExamsScoresOutputDTO, error) {
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

	// [STEP 2] Verify user is teacher or is requested user
	if user.Role != user_constants.UserRoleProfessor && user.ID != input.UserID {
		return nil, fmt.Errorf("User %q does not have permissions to view this exam scoring", user.Username)
	}

	// [STEP 3] Get exams scores by userID
	examScores, err := uc.examScoreRepository.GetExamScores(ctx, nil, &input.UserID)
	if err != nil {
		return nil, err
	}

	if len(examScores) == 0 {
		return nil, fmt.Errorf("no exam scores found for userID %q", input.UserID)
	}

	// [STEP 4] Aggregate scores by exam
	aggregatedScores := services.AggregateScoresByExam(examScores)

	// [STEP 4] If user is teacher, filter examScores, only exams owned by the teacher
	if user.Role == user_constants.UserRoleProfessor {
		// [STEP 4.1] Get all exams from examScores
		exams, err := uc.getExamsFromExamScores(ctx, aggregatedScores)
		if err != nil {
			return nil, err
		}

		// [STEP 3.2] Filter examScores, only exams owned by the teacher
		filteredExamScores := services.FilterExamScoresByExams(aggregatedScores, exams)
		aggregatedScores = filteredExamScores
	}

	// [STEP 5] Map list of examScores to examScoringDTO
	var dtos []*dtos.ExamScoringDTO
	for _, scores := range aggregatedScores {
		best_score := services.CalculateBestScore(scores)
		dto, err := mapper.MapExamScoreEntitiesToExamScoringDTO(scores, best_score)
		if err != nil {
			return nil, err
		}

		dtos = append(dtos, dto)
	}

	// [STEP 6] Calculate prom score
	var prom_score float64
	for _, scores := range aggregatedScores {
		prom_score += services.CalculatePromScore(scores)
	}

	// [STEP 7] Map to output DTO and return
	result, err := mapper.MapExamScoresDTOtoUserExamScoresDTO(dtos, prom_score, input.UserID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (uc *GetExamsScoresByUserUseCase) getExamsFromExamScores(ctx context.Context, aggregatedScores map[string][]*examEntities.ExamScore) ([]examEntities.Exam, error) {
	var exams []examEntities.Exam
	for exam := range aggregatedScores {
		exam, err := uc.examRepository.GetExamByID(ctx, exam)
		if err != nil {
			return nil, err
		}

		if exam == nil {
			continue
		}

		exams = append(exams, *exam)
	}

	return exams, nil
}



