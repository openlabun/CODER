package roble_infrastructure

import (
	"context"
	"fmt"
	"strings"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
)

type SubmissionResultRepository struct {
	adapter *infrastructure.RobleDatabaseAdapter
}

func NewSubmissionResultRepository(adapter *infrastructure.RobleDatabaseAdapter) *SubmissionResultRepository {
	return &SubmissionResultRepository{adapter: adapter}
}

func (r *SubmissionResultRepository) CreateResult(ctx context.Context, result *Entities.SubmissionResult) (*Entities.SubmissionResult, error) {
	if result == nil {
		return nil, fmt.Errorf("result is nil")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	_, err := r.adapter.Insert(submissionResultTableName, []map[string]any{resultToRecord(result)})
	if err != nil {
		return nil, err
	}

	return r.GetResultByID(ctx, result.ID)
}

func (r *SubmissionResultRepository) UpdateResult(ctx context.Context, result *Entities.SubmissionResult) (*Entities.SubmissionResult, error) {
	if result == nil {
		return nil, fmt.Errorf("result is nil")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	resultID := strings.TrimSpace(result.ID)
	if resultID == "" {
		return nil, fmt.Errorf("result id is required")
	}

	_, err := r.adapter.Update(submissionResultTableName, "ID", resultID, resultToUpdates(result))
	if err != nil {
		return nil, err
	}

	return r.GetResultByID(ctx, resultID)
}

func (r *SubmissionResultRepository) DeleteResult(ctx context.Context, resultID string) error {
	normalizedID := strings.TrimSpace(resultID)
	if normalizedID == "" {
		return fmt.Errorf("resultID is required")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return err
	}

	_, err := r.adapter.Delete(submissionResultTableName, "ID", normalizedID)
	return err
}

func (r *SubmissionResultRepository) GetResultByID(ctx context.Context, resultID string) (*Entities.SubmissionResult, error) {
	normalizedID := strings.TrimSpace(resultID)
	if normalizedID == "" {
		return nil, fmt.Errorf("resultID is required")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	res, err := r.adapter.Read(submissionResultTableName, map[string]string{"ID": normalizedID})
	if err != nil {
		return nil, err
	}

	record, err := firstRecord(res)
	if err != nil {
		return nil, nil
	}

	actualOutput, err := getIOVariableByID(ctx, r.adapter, asString(record["ActualOutput"]))
	if err != nil {
		return nil, err
	}

	return recordToResult(record, actualOutput)
}

func (r *SubmissionResultRepository) GetResultsBySubmissionID(ctx context.Context, submissionID string) ([]*Entities.SubmissionResult, error) {
	return r.getResultsByField(ctx, "SubmissionID", submissionID)
}

func (r *SubmissionResultRepository) GetResultByTestCase(ctx context.Context, testCaseID string) ([]*Entities.SubmissionResult, error) {
	return r.getResultsByField(ctx, "TestCaseID", testCaseID)
}

func (r *SubmissionResultRepository) getResultsByField(ctx context.Context, field, value string) ([]*Entities.SubmissionResult, error) {
	normalizedValue := strings.TrimSpace(value)
	if normalizedValue == "" {
		return nil, fmt.Errorf("%s is required", field)
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	res, err := r.adapter.Read(submissionResultTableName, map[string]string{field: normalizedValue})
	if err != nil {
		return nil, err
	}

	records := extractRecords(res)
	if len(records) == 0 {
		return []*Entities.SubmissionResult{}, nil
	}

	results := make([]*Entities.SubmissionResult, 0, len(records))
	for _, record := range records {
		actualOutput, fetchErr := getIOVariableByID(ctx, r.adapter, asString(record["ActualOutput"]))
		if fetchErr != nil {
			return nil, fetchErr
		}

		result, mapErr := recordToResult(record, actualOutput)
		if mapErr != nil {
			return nil, mapErr
		}
		if result != nil {
			results = append(results, result)
		}
	}

	return results, nil
}
