package submission_usecases

import (
	"fmt"
	"context"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
	mapper "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission/mapper"

	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	user_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/user"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	submissionRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
)

type GetItemScoringUseCase struct {
	userRepository userRepository.UserRepository
	examItemScoreRepository examRepository.ExamItemScoreRepository
	submissionRepository submissionRepository.SubmissionRepository
	examItemRepository examRepository.ExamItemRepository
	examScoreRepository examRepository.ExamScoreRepository
}

func NewGetItemScoringUseCase(userRepository userRepository.UserRepository, examItemScoreRepository examRepository.ExamItemScoreRepository, submissionRepository submissionRepository.SubmissionRepository, examItemRepository examRepository.ExamItemRepository, examScoreRepository examRepository.ExamScoreRepository) *GetItemScoringUseCase {
	return &GetItemScoringUseCase{
		userRepository: userRepository,
		examItemScoreRepository: examItemScoreRepository,
		submissionRepository: submissionRepository,
		examItemRepository: examItemRepository,
		examScoreRepository: examScoreRepository,
	}
}

func (uc *GetItemScoringUseCase) Execute(ctx context.Context, input dtos.GetItemScoringInput) (*dtos.ExamScoreDetailsOutputDTO, error) {
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

	// [STEP 2] Retrieve exam score and validate it exists
	examScore, err := uc.examScoreRepository.GetExamScoreByID(ctx, input.ExamScoreID)
	if err != nil {
		return nil, err
	}

	if examScore == nil {
		return nil, fmt.Errorf("ExamScore with ID %q does not exist", input.ExamScoreID)
	}

	// [STEP 3] Verify user is teacher or student owner of the exam score
	if user.Role != user_constants.UserRoleProfessor && user.ID != examScore.StudentID {
		return nil, fmt.Errorf("User %q does not have permissions to view this exam scoring", user.Username)
	}

	// [STEP 4] Get exam items scores for the given exam score ID
	examItemScores, err := uc.examItemScoreRepository.GetExamItemScoresByExamScoreID(ctx, input.ExamScoreID)
	if err != nil {
		return nil, err
	}

	// [STEP 5] For each exam item score, retrieve submissions and map to ExamItemScoringDTO
	var dtos []*dtos.ExamItemScoringDTO
	for _, examItemScore := range examItemScores {

		// [STEP 5.1] Get exam item for its exam item score
		examItem, err := uc.examItemRepository.GetExamItemByID(ctx, examItemScore.ExamItemID)
		if err != nil {
			return nil, err
		}

		// [STEP 5.2] Get submissions by exam item score ID
		submissions, err := uc.submissionRepository.GetSubmissionsByExamItemScoreID(ctx, examItem.ID)
		if err != nil {
			return nil, err
		}

		// [STEP 5.3] Map exam item score and submissions to ExamItemScoringDTO
		examItemScoringDTO, err := mapper.MapExamItemScoreAndSubmissionsToExamItemScoringDTO(examItem, examItemScore, submissions)
		if err != nil {
			return nil, err
		}

		dtos = append(dtos, examItemScoringDTO)
	}

	// [STEP 6] Map to output DTO and return
	examScoreDetailsDTO, err := mapper.MapExamItemScoringDTOToExamScoreDetailsOutputDTO(dtos, input.ExamScoreID, examScore.ExamID)
	if err != nil {
		return nil, err
	}

	return examScoreDetailsDTO, nil
}