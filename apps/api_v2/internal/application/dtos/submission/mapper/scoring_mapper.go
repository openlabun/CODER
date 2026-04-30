package mapper

import (
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	SubmissionEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
)

func MapExamScoreEntitiesToExamScoringDTO(entities []*Entities.ExamScore, bestScore int) (*dtos.ExamScoringDTO, error) {
	if len(entities) == 0 {
		return nil, fmt.Errorf("Entities list cannot be empty")
	}

	var entities_list []Entities.ExamScore
	for _, entity := range entities {
		if entity != nil {
			entities_list = append(entities_list, *entity)
		}
	}

	return &dtos.ExamScoringDTO{
		ExamID:    entities[0].ExamID,
		BestScore: bestScore,
		ExamScores: entities_list,
	}, nil
}

func MapExamScoresDTOtoUserExamScoresDTO (examScoresDTO []*dtos.ExamScoringDTO, prom_score float64, userID string) (*dtos.UserExamsScoresOutputDTO, error) {
	if len(examScoresDTO) == 0 {
		return nil, fmt.Errorf("ExamScoresDTO list cannot be empty")
	}

	examScores := make([]dtos.ExamScoringDTO, len(examScoresDTO))
	for i, examScoreDTO := range examScoresDTO {
		if examScoreDTO == nil {
			continue
		}
		examScores[i] = *examScoreDTO
	}

	return &dtos.UserExamsScoresOutputDTO{
		UserID: userID,
		PromScore: prom_score,
		ExamScores: examScores,
	}, nil
}

func MapExamScoresDTOtoUserExamScoresOutputDTO (examScoresDTO *dtos.ExamScoringDTO, userID string) (*dtos.UserExamScoresOutputDTO, error) {
	if examScoresDTO == nil {
		return nil, fmt.Errorf("ExamScoresDTO cannot be nil")
	}

	return &dtos.UserExamScoresOutputDTO{
		UserID: userID,
		ExamScores: *examScoresDTO,
	}, nil
}

func MapUserExamScoresOutputDTOToUsersExamScoresOutputDTO(userExamScoresDTOs []*dtos.UserExamScoresOutputDTO, examID string) (*dtos.UsersExamScoresOutputDTO, error) {
	if len(userExamScoresDTOs) == 0 {
		return nil, fmt.Errorf("UserExamScoresDTOs list cannot be empty")
	}

	userExamScores := make([]dtos.UserExamScoresOutputDTO, len(userExamScoresDTOs))
	for i, userExamScoresDTO := range userExamScoresDTOs {
		if userExamScoresDTO == nil {
			continue
		}

		userExamScores[i] = *userExamScoresDTO
	}		

	return &dtos.UsersExamScoresOutputDTO{
		ExamID: examID,
		ExamScores: userExamScores,
	}, nil
}

func MapExamItemScoreAndSubmissionsToExamItemScoringDTO(examItem *Entities.ExamItem, examItemScore *Entities.ExamItemScore, submissions []*SubmissionEntities.Submission) (*dtos.ExamItemScoringDTO, error) {
	if examItem == nil {
		return nil, fmt.Errorf("ExamItem cannot be nil")
	}

	if examItemScore == nil {
		return nil, fmt.Errorf("ExamItemScore cannot be nil")
	}

	submissionsList := make([]SubmissionEntities.Submission, len(submissions))
	for i, submission := range submissions {
		if submission != nil {
			submissionsList[i] = *submission
		}
	}

	return &dtos.ExamItemScoringDTO{
		ExamItem: *examItem,
		ExamItemScore: *examItemScore,
		Submissions: submissionsList,
	}, nil
}

func MapExamItemScoringDTOToExamScoreDetailsOutputDTO (examItemScoringDTOs []*dtos.ExamItemScoringDTO, examScoreID string, examID string) (*dtos.ExamScoreDetailsOutputDTO, error) {
	if len(examItemScoringDTOs) == 0 {
		return nil, fmt.Errorf("ExamItemScoringDTOs list cannot be empty")
	}

	examItems := make([]dtos.ExamItemScoringDTO, len(examItemScoringDTOs))
	for i, examItemScoringDTO := range examItemScoringDTOs {
		if examItemScoringDTO == nil {
			continue
		}

		examItems[i] = *examItemScoringDTO
	}

	return &dtos.ExamScoreDetailsOutputDTO{
		ExamScoreID: examScoreID,
		ExamID: examID,
		ExamItems: examItems,
	}, nil
}

