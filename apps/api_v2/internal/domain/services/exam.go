package services

import (
	"context"

	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
)

func RemoveExam (ctx context.Context,
	examID string,
	examRepository examRepository.ExamRepository,
	examItemRepository examRepository.ExamItemRepository,
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

		// [STEP 3] Delete the exam itself
		err = examRepository.DeleteExam(ctx, examID)
		if err != nil {
			return err
		}

		return nil
	}