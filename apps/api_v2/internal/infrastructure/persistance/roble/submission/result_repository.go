package roble_infrastructure

import (
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

func (r *SubmissionResultRepository) CreateResult(result *Entities.SubmissionResult) (*Entities.SubmissionResult, error) {
	if result == nil {
		return nil, fmt.Errorf("result is nil")
	}

	if result.ActualOutput != nil {
		if err := upsertIOVariable(r.adapter, *result.ActualOutput); err != nil {
			return nil, err
		}
	}

	_, err := r.adapter.Insert(submissionResultTableName, []map[string]any{resultToRecord(result)})
	if err != nil {
		return nil, err
	}

	return r.GetResultByID(result.ID)
}

func (r *SubmissionResultRepository) UpdateResult(result *Entities.SubmissionResult) (*Entities.SubmissionResult, error) {
	if result == nil {
		return nil, fmt.Errorf("result is nil")
	}

	resultID := strings.TrimSpace(result.ID)
	if resultID == "" {
		return nil, fmt.Errorf("result id is required")
	}

	if result.ActualOutput != nil {
		if err := upsertIOVariable(r.adapter, *result.ActualOutput); err != nil {
			return nil, err
		}
	}

	_, err := r.adapter.Update(submissionResultTableName, "ID", resultID, resultToUpdates(result))
	if err != nil {
		return nil, err
	}

	return r.GetResultByID(resultID)
}

func (r *SubmissionResultRepository) DeleteResult(resultID string) error {
	normalizedID := strings.TrimSpace(resultID)
	if normalizedID == "" {
		return fmt.Errorf("resultID is required")
	}

	res, err := r.adapter.Read(submissionResultTableName, map[string]string{"ID": normalizedID})
	if err != nil {
		return err
	}

	var ActualOutput string
	if record, findErr := firstRecord(res); findErr == nil {
		ActualOutput = strings.TrimSpace(asString(record["ActualOutput"]))
	}

	_, err = r.adapter.Delete(submissionResultTableName, "ID", normalizedID)
	if err != nil {
		return err
	}

	return deleteIOVariableByID(r.adapter, ActualOutput)
}

func (r *SubmissionResultRepository) GetResultByID(resultID string) (*Entities.SubmissionResult, error) {
	normalizedID := strings.TrimSpace(resultID)
	if normalizedID == "" {
		return nil, fmt.Errorf("resultID is required")
	}

	res, err := r.adapter.Read(submissionResultTableName, map[string]string{"ID": normalizedID})
	if err != nil {
		return nil, err
	}

	record, err := firstRecord(res)
	if err != nil {
		return nil, nil
	}

	actualOutput, err := getIOVariableByID(r.adapter, asString(record["ActualOutput"]))
	if err != nil {
		return nil, err
	}

	return recordToResult(record, actualOutput)
}

func (r *SubmissionResultRepository) GetResultsBySubmissionID(submissionID string) ([]*Entities.SubmissionResult, error) {
	return r.getResultsByField("SubmissionID", submissionID)
}

func (r *SubmissionResultRepository) GetResultByTestCase(testCaseID string) ([]*Entities.SubmissionResult, error) {
	return r.getResultsByField("TestCaseID", testCaseID)
}

func (r *SubmissionResultRepository) getResultsByField(field, value string) ([]*Entities.SubmissionResult, error) {
	normalizedValue := strings.TrimSpace(value)
	if normalizedValue == "" {
		return nil, fmt.Errorf("%s is required", field)
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
		actualOutput, fetchErr := getIOVariableByID(r.adapter, asString(record["ActualOutput"]))
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
