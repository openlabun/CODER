package services

import (
	"context"

	examEntity "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	submissionRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
)

func CalculateExamScore(ctx context.Context, examScore *examEntity.ExamScore, examScoreRepository examRepository.ExamScoreRepository, examItemRepository examRepository.ExamItemRepository, examItemScoreRepository examRepository.ExamItemScoreRepository, submissionRepository submissionRepository.SubmissionRepository) (*examEntity.ExamScore, error) {
	// [STEP 1] Get all ExamItemScores for the ExamScore
	examItemScores, err := examItemScoreRepository.GetExamItemScoresByExamScoreID(ctx, examScore.ID)
	if err != nil {
		return nil, err
	}

	// [STEP 2] For each ExamItemScore, get best submission and update score
	totalScore := 0
	for _, itemScore := range examItemScores {
		// [STEP 2.1] Get best submission for the ExamItemScore
		submission, err := submissionRepository.GetBestSubmissionByExamItemScoreID(ctx, itemScore.ID)
		if err != nil {
			return nil, err
		}

		if submission == nil {
			continue
		}

		// [STEP 2.2] Update ExamItemScore with the new score
		itemScore.Score = submission.Score
		_, err = examItemScoreRepository.UpdateExamItemScore(ctx, itemScore)
		if err != nil {
			return nil, err
		}

		// [STEP 2.3] Get ExamItem for the ExamItemScore
		examItem, err := examItemRepository.GetExamItemByID(ctx, itemScore.ExamItemID)
		if err != nil {
			return nil, err
		}

		if examItem == nil {
			continue
		}

		// [STEP 2.4] Add score to total score
		totalScore += itemScore.Score * examItem.Points

	}

	// [STEP 3] Update ExamScore with the new total score
	examScore.Score = totalScore/100
	updatedExamScore, err := examScoreRepository.UpdateExamScore(ctx, examScore)
	if err != nil {
		return nil, err
	}

	return updatedExamScore, nil
}

func RemoveExamScore(ctx context.Context,
	examScoreID string,
	examScoreRepository examRepository.ExamScoreRepository,
	examItemScoreRepository examRepository.ExamItemScoreRepository,
) error {
	// [STEP 1] Get ExamScore by ID
	examScore, err := examScoreRepository.GetExamScoreByID(ctx, examScoreID)
	if err != nil {
		return err
	}

	if examScore == nil {
		return nil
	}

	// [STEP 2] Get all ExamItemScores for the ExamScore
	examItemScores, err := examItemScoreRepository.GetExamItemScoresByExamScoreID(ctx, examScoreID)
	if err != nil {
		return err
	}

	// [STEP 3] Delete all ExamItemScores for the ExamScore
	for _, itemScore := range examItemScores {
		if itemScore != nil {
			err = examItemScoreRepository.DeleteExamItemScore(ctx, itemScore.ID)
			if err != nil {
				return err
			}
		}
	}

	// [STEP 4] Delete the ExamScore itself
	err = examScoreRepository.DeleteExamScore(ctx, examScoreID)
	if err != nil {
		return err
	}

	return nil
}