package services

import (
	"context"

	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
)

func RemoveExam (ctx context.Context,
	examID string,
	examRepository examRepository.ExamRepository,
	examItemRepository examRepository.ExamItemRepository,
	examScoreRepository examRepository.ExamScoreRepository,
	examItemScoreRepository examRepository.ExamItemScoreRepository,
	) error {
		// [STEP 1] Get all exam items for the exam
		examItems, err := examItemRepository.GetExamItem(ctx, &examID, nil)
		if err != nil {
			return err
		}

		// [STEP 2] Delete all existing exam items for the exam
		for _, item := range examItems {
			if item != nil {
				err = examItemRepository.DeleteExamItem(ctx, item.ID)
				if err != nil {
					return err
				}
			}
		}

		// [STEP 3] Get all exam scores for the exam
		examScores, err := examScoreRepository.GetExamScores(ctx, &examID, nil)
		if err != nil {
			return err
		}

		// [STEP 4] Delete all existing exam scores for the exam
		for _, score := range examScores {
			if score != nil {
				err = RemoveExamScore(ctx, score.ID, examScoreRepository, examItemScoreRepository)
				if err != nil {
					return err
				}
			}
		}

		// [STEP 5] Delete the exam itself
		err = examRepository.DeleteExam(ctx, examID)
		if err != nil {
			return err
		}

		return nil
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