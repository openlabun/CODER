package roble_infrastructure

import (
	"context"
	"fmt"
	"strings"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
)

const examItemTableName = "ExamItem"

type ExamItemRepository struct {
	adapter *infrastructure.RobleDatabaseAdapter
}

func NewExamItemRepository(adapter *infrastructure.RobleDatabaseAdapter) *ExamItemRepository {
	return &ExamItemRepository{adapter: adapter}
}

func (r *ExamItemRepository) CreateExamItem(ctx context.Context, item *Entities.ExamItem) (*Entities.ExamItem, error) {
	if item == nil {
		return nil, fmt.Errorf("exam item is nil")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}
	_, err := r.adapter.Insert(examItemTableName, []map[string]any{examItemToRecord(item)})
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (r *ExamItemRepository) UpdateExamItem(ctx context.Context, item *Entities.ExamItem) (*Entities.ExamItem, error) {
	if item == nil {
		return nil, fmt.Errorf("exam item is nil")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}
	itemID := strings.TrimSpace(item.ID)
	if itemID == "" {
		return nil, fmt.Errorf("exam item id is required")
	}
	updates := examItemToUpdates(item)
	_, err := r.adapter.Update(examItemTableName, "ID", itemID, updates)
	if err != nil {
		return nil, err
	}
	item.ID = itemID
	return item, nil
}

func (r *ExamItemRepository) DeleteExamItem(ctx context.Context, examItemID string) error {
	normalizedID := strings.TrimSpace(examItemID)
	if normalizedID == "" {
		return fmt.Errorf("examItemID is required")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return err
	}
	_, err := r.adapter.Delete(examItemTableName, "ID", normalizedID)
	return err
}

func (r *ExamItemRepository) GetExamItemByID(ctx context.Context, examItemID string) (*Entities.ExamItem, error) {
	normalizedID := strings.TrimSpace(examItemID)
	if normalizedID == "" {
		return nil, fmt.Errorf("examItemID is required")
	}

	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	res, err := r.adapter.Read(examItemTableName, map[string]string{"ID": normalizedID})
	if err != nil {
		return nil, err
	}

	records := extractRecords(res)
	if len(records) == 0 {
		return nil, nil // No record found, return nil without error
	}

	return recordToExamItem(records[0])
}

func (r *ExamItemRepository) GetExamItem(ctx context.Context, examID *string, challengeID *string) ([]*Entities.ExamItem, error) {
	filters := map[string]string{}
	if examID != nil && strings.TrimSpace(*examID) != "" {
		filters["ExamID"] = strings.TrimSpace(*examID)
	}
	if challengeID != nil && strings.TrimSpace(*challengeID) != "" {
		filters["ChallengeID"] = strings.TrimSpace(*challengeID)
	}
	if len(filters) == 0 {
		return nil, fmt.Errorf("at least one of examID or challengeID must be provided")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}
	res, err := r.adapter.Read(examItemTableName, filters)
	if err != nil {
		return nil, err
	}
	records := extractRecords(res)
	if len(records) == 0 {
		return []*Entities.ExamItem{}, nil
	}
	items := make([]*Entities.ExamItem, 0, len(records))
	for _, record := range records {
		item, mapErr := recordToExamItem(record)
		if mapErr != nil {
			return nil, mapErr
		}
		if item != nil {
			items = append(items, item)
		}
	}
	return items, nil
}
