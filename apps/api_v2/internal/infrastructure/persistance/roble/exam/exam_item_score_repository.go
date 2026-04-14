package roble_infrastructure

import (
	"context"
	"fmt"
	"strings"
	"time"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
)

const examItemScoreTableName = "ExamItemScore"

type ExamItemScoreRepository struct {
	adapter *infrastructure.RobleDatabaseAdapter
}

func NewExamItemScoreRepository(adapter *infrastructure.RobleDatabaseAdapter) *ExamItemScoreRepository {
	return &ExamItemScoreRepository{adapter: adapter}
}

func (r *ExamItemScoreRepository) CreateExamItemScore(ctx context.Context, examItemScore *Entities.ExamItemScore) (*Entities.ExamItemScore, error) {
	if examItemScore == nil {
		return nil, fmt.Errorf("exam item score is nil")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	_, err := r.adapter.Insert(examItemScoreTableName, []map[string]any{examItemScoreToRecord(examItemScore)})
	if err != nil {
		return nil, err
	}

	return examItemScore, nil
}

func (r *ExamItemScoreRepository) UpdateExamItemScore(ctx context.Context, examItemScore *Entities.ExamItemScore) (*Entities.ExamItemScore, error) {
	if examItemScore == nil {
		return nil, fmt.Errorf("exam item score is nil")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	examItemScoreID := strings.TrimSpace(examItemScore.ID)
	if examItemScoreID == "" {
		return nil, fmt.Errorf("exam item score id is required")
	}

	updates := examItemScoreToUpdates(examItemScore)
	updates["UpdatedAt"] = time.Now().UTC().Format(time.RFC3339)

	_, err := r.adapter.Update(examItemScoreTableName, "ID", examItemScoreID, updates)
	if err != nil {
		return nil, err
	}

	examItemScore.ID = examItemScoreID
	if ts, ok := updates["UpdatedAt"].(string); ok {
		examItemScore.UpdatedAt = ts
	}

	return examItemScore, nil
}

func (r *ExamItemScoreRepository) DeleteExamItemScore(ctx context.Context, examItemScoreID string) error {
	normalizedID := strings.TrimSpace(examItemScoreID)
	if normalizedID == "" {
		return fmt.Errorf("examItemScoreID is required")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return err
	}

	_, err := r.adapter.Delete(examItemScoreTableName, "ID", normalizedID)
	return err
}

func (r *ExamItemScoreRepository) GetExamItemScoreByID(ctx context.Context, examItemScoreID string) (*Entities.ExamItemScore, error) {
	normalizedID := strings.TrimSpace(examItemScoreID)
	if normalizedID == "" {
		return nil, fmt.Errorf("examItemScoreID is required")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	res, err := r.adapter.Read(examItemScoreTableName, map[string]string{"ID": normalizedID})
	if err != nil {
		return nil, err
	}

	record, err := firstRecord(res)
	if err != nil {
		return nil, nil
	}

	return recordToExamItemScore(record)
}

func (r *ExamItemScoreRepository) GetExamItemScoresByExamScoreID(ctx context.Context, examScoreID string) ([]*Entities.ExamItemScore, error) {
	normalizedExamScoreID := strings.TrimSpace(examScoreID)
	if normalizedExamScoreID == "" {
		return nil, fmt.Errorf("examScoreID is required")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	res, err := r.adapter.Read(examItemScoreTableName, map[string]string{"ExamScoreID": normalizedExamScoreID})
	if err != nil {
		return nil, err
	}

	records := extractRecords(res)
	if len(records) == 0 {
		return []*Entities.ExamItemScore{}, nil
	}

	examItemScores := make([]*Entities.ExamItemScore, 0, len(records))
	for _, record := range records {
		examItemScore, mapErr := recordToExamItemScore(record)
		if mapErr != nil {
			return nil, mapErr
		}
		if examItemScore != nil {
			examItemScores = append(examItemScores, examItemScore)
		}
	}

	return examItemScores, nil
}

func (r *ExamItemScoreRepository) GetExamItemScore(ctx context.Context, examItemID, examScoreID string) (*Entities.ExamItemScore, error) {
	normalizedExamItemID := strings.TrimSpace(examItemID)
	normalizedExamScoreID := strings.TrimSpace(examScoreID)

	if normalizedExamItemID == "" {
		return nil, fmt.Errorf("examItemID is required")
	}
	if normalizedExamScoreID == "" {
		return nil, fmt.Errorf("examScoreID is required")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	res, err := r.adapter.Read(examItemScoreTableName, map[string]string{
		"ExamItemID":  normalizedExamItemID,
		"ExamScoreID": normalizedExamScoreID,
	})
	if err != nil {
		return nil, err
	}

	record, err := firstRecord(res)
	if err != nil {
		return nil, nil
	}

	return recordToExamItemScore(record)
}
