package roble_infrastructure

import (
	"context"
	"fmt"
	"strings"
	"time"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
)

const examScoreTableName = "ExamScore"

type ExamScoreRepository struct {
	adapter *infrastructure.RobleDatabaseAdapter
}

func NewExamScoreRepository(adapter *infrastructure.RobleDatabaseAdapter) *ExamScoreRepository {
	return &ExamScoreRepository{adapter: adapter}
}

func (r *ExamScoreRepository) CreateExamScore(ctx context.Context, examScore *Entities.ExamScore) (*Entities.ExamScore, error) {
	if examScore == nil {
		return nil, fmt.Errorf("exam score is nil")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	_, err := r.adapter.Insert(examScoreTableName, []map[string]any{examScoreToRecord(examScore)})
	if err != nil {
		return nil, err
	}

	return examScore, nil
}

func (r *ExamScoreRepository) UpdateExamScore(ctx context.Context, examScore *Entities.ExamScore) (*Entities.ExamScore, error) {
	if examScore == nil {
		return nil, fmt.Errorf("exam score is nil")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	examScoreID := strings.TrimSpace(examScore.ID)
	if examScoreID == "" {
		return nil, fmt.Errorf("exam score id is required")
	}

	updates := examScoreToUpdates(examScore)
	updates["UpdatedAt"] = time.Now().UTC().Format(time.RFC3339)

	_, err := r.adapter.Update(examScoreTableName, "ID", examScoreID, updates)
	if err != nil {
		return nil, err
	}

	examScore.ID = examScoreID
	if ts, ok := updates["UpdatedAt"].(string); ok {
		examScore.UpdatedAt = ts
	}

	return examScore, nil
}

func (r *ExamScoreRepository) DeleteExamScore(ctx context.Context, examScoreID string) error {
	normalizedID := strings.TrimSpace(examScoreID)
	if normalizedID == "" {
		return fmt.Errorf("examScoreID is required")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return err
	}

	_, err := r.adapter.Delete(examScoreTableName, "ID", normalizedID)
	return err
}

func (r *ExamScoreRepository) GetExamScores(ctx context.Context, examID, studentID *string) ([]*Entities.ExamScore, error) {
	filters := map[string]string{}

	if examID != nil {
		if normalizedExamID := strings.TrimSpace(*examID); normalizedExamID != "" {
			filters["ExamID"] = normalizedExamID
		}
	}

	if studentID != nil {
		if normalizedStudentID := strings.TrimSpace(*studentID); normalizedStudentID != "" {
			filters["StudentID"] = normalizedStudentID
		}
	}

	if len(filters) == 0 {
		return nil, fmt.Errorf("at least one of examID or studentID must be provided")
	}

	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	res, err := r.adapter.Read(examScoreTableName, filters)
	if err != nil {
		return nil, err
	}

	records := extractRecords(res)
	if len(records) == 0 {
		return []*Entities.ExamScore{}, nil
	}

	examScores := make([]*Entities.ExamScore, 0, len(records))
	for _, record := range records {
		examScore, mapErr := recordToExamScore(record)
		if mapErr != nil {
			return nil, mapErr
		}
		if examScore != nil {
			examScores = append(examScores, examScore)
		}
	}

	return examScores, nil
}

func (r *ExamScoreRepository) GetExamScoreByID(ctx context.Context, examScoreID string) (*Entities.ExamScore, error) {
	normalizedID := strings.TrimSpace(examScoreID)
	if normalizedID == "" {
		return nil, fmt.Errorf("examScoreID is required")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	res, err := r.adapter.Read(examScoreTableName, map[string]string{"ID": normalizedID})
	if err != nil {
		return nil, err
	}

	record, err := firstRecord(res)
	if err != nil {
		return nil, nil
	}

	return recordToExamScore(record)
}

func (r *ExamScoreRepository) GetExamScoreBySessionID(ctx context.Context, sessionID string) (*Entities.ExamScore, error) {
	normalizedSessionID := strings.TrimSpace(sessionID)
	if normalizedSessionID == "" {
		return nil, fmt.Errorf("sessionID is required")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	res, err := r.adapter.Read(examScoreTableName, map[string]string{"SessionID": normalizedSessionID})
	if err != nil {
		return nil, err
	}

	record, err := firstRecord(res)
	if err != nil {
		return nil, nil
	}

	return recordToExamScore(record)
}
