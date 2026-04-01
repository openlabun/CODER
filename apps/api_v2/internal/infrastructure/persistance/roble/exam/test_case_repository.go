package roble_infrastructure

import (
	"context"
	"fmt"
	"strings"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
)

const (
	testCaseTableName = "TestCase"
)

type TestCaseRepository struct {
	adapter *infrastructure.RobleDatabaseAdapter
}

func NewTestCaseRepository(adapter *infrastructure.RobleDatabaseAdapter) *TestCaseRepository {
	return &TestCaseRepository{adapter: adapter}
}

func (r *TestCaseRepository) CreateTestCase(ctx context.Context, testCase *Entities.TestCase) (*Entities.TestCase, error) {
	if testCase == nil {
		return nil, fmt.Errorf("testCase is nil")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	_, err := r.adapter.Insert(testCaseTableName, []map[string]any{testCaseToRecord(testCase)})
	if err != nil {
		return nil, err
	}

	return r.GetTestCaseByID(ctx, testCase.ID)
}

func (r *TestCaseRepository) UpdateTestCase(ctx context.Context, testCase *Entities.TestCase) (*Entities.TestCase, error) {
	if testCase == nil {
		return nil, fmt.Errorf("testCase is nil")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	testCaseID := strings.TrimSpace(testCase.ID)
	if testCaseID == "" {
		return nil, fmt.Errorf("testCase id is required")
	}

	_, err := r.adapter.Update(testCaseTableName, "ID", testCaseID, testCaseToUpdates(testCase))
	if err != nil {
		return nil, err
	}

	return r.GetTestCaseByID(ctx, testCaseID)
}

func (r *TestCaseRepository) DeleteTestCase(ctx context.Context, testCaseID string) error {
	normalizedID := strings.TrimSpace(testCaseID)
	if normalizedID == "" {
		return fmt.Errorf("testCaseID is required")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return err
	}

	_, err := r.adapter.Delete(testCaseTableName, "ID", normalizedID)
	return err
}

func (r *TestCaseRepository) GetTestCaseByID(ctx context.Context, testCaseID string) (*Entities.TestCase, error) {
	normalizedID := strings.TrimSpace(testCaseID)
	if normalizedID == "" {
		return nil, fmt.Errorf("testCaseID is required")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	res, err := r.adapter.Read(testCaseTableName, map[string]string{"ID": normalizedID})
	if err != nil {
		return nil, err
	}

	record, err := firstRecord(res)
	if err != nil {
		return nil, nil
	}

	return r.recordToHydratedTestCase(ctx, record)
}

func (r *TestCaseRepository) GetTestCasesByChallengeID(ctx context.Context, challengeID string) ([]*Entities.TestCase, error) {
	normalizedID := strings.TrimSpace(challengeID)
	if normalizedID == "" {
		return nil, fmt.Errorf("challengeID is required")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	res, err := r.adapter.Read(testCaseTableName, map[string]string{"ChallengeID": normalizedID})
	if err != nil {
		return nil, err
	}

	records := extractRecords(res)
	if len(records) == 0 {
		return []*Entities.TestCase{}, nil
	}

	testCases := make([]*Entities.TestCase, 0, len(records))
	for _, record := range records {
		testCase, mapErr := r.recordToHydratedTestCase(ctx, record)
		if mapErr != nil {
			return nil, mapErr
		}
		if testCase != nil {
			testCases = append(testCases, testCase)
		}
	}

	return testCases, nil
}

func (r *TestCaseRepository) GetInputVariablesByTestCaseID(ctx context.Context, testCaseID string) ([]*Entities.IOVariable, error) {
	testCase, err := r.GetTestCaseByID(ctx, testCaseID)
	if err != nil {
		return nil, err
	}
	if testCase == nil {
		return []*Entities.IOVariable{}, nil
	}

	variables := make([]*Entities.IOVariable, 0, len(testCase.Input))
	for i := range testCase.Input {
		variable := testCase.Input[i]
		variables = append(variables, &variable)
	}

	return variables, nil
}

func (r *TestCaseRepository) GetOutputVariablesByTestCaseID(ctx context.Context, testCaseID string) ([]*Entities.IOVariable, error) {
	testCase, err := r.GetTestCaseByID(ctx, testCaseID)
	if err != nil {
		return nil, err
	}
	if testCase == nil || strings.TrimSpace(testCase.ExpectedOutput.ID) == "" {
		return []*Entities.IOVariable{}, nil
	}

	output := testCase.ExpectedOutput
	return []*Entities.IOVariable{&output}, nil
}

func (r *TestCaseRepository) recordToHydratedTestCase(ctx context.Context, record map[string]any) (*Entities.TestCase, error) {
	inputIDs := asStringList(record["Input"])
	if len(inputIDs) == 0 {
		inputIDs = asStringList(record["InputVariables"])
	}

	outputID := asString(record["ExpectedOutput"])
	if strings.TrimSpace(outputID) == "" {
		outputID = asString(record["OutputVariable"])
	}

	inputVariables, err := r.getIOVariablesByIDs(ctx, inputIDs)
	if err != nil {
		return nil, err
	}

	outputVariable, err := r.getIOVariableByID(ctx, outputID)
	if err != nil {
		return nil, err
	}

	return recordToTestCase(record, inputVariables, outputVariable)
}

func (r *TestCaseRepository) upsertTestCaseIOVariables(ctx context.Context, testCase *Entities.TestCase) error {
	for _, input := range testCase.Input {
		if err := r.upsertIOVariable(ctx, input); err != nil {
			return err
		}
	}

	if err := r.upsertIOVariable(ctx, testCase.ExpectedOutput); err != nil {
		return err
	}

	return nil
}

func (r *TestCaseRepository) getIOVariablesByIDs(ctx context.Context, ids []string) ([]Entities.IOVariable, error) {
	if len(ids) == 0 {
		return []Entities.IOVariable{}, nil
	}

	variables := make([]Entities.IOVariable, 0, len(ids))
	for _, id := range ids {
		ioVariable, err := r.getIOVariableByID(ctx, id)
		if err != nil {
			return nil, err
		}
		if ioVariable != nil {
			variables = append(variables, *ioVariable)
		}
	}

	return variables, nil
}

func (r *TestCaseRepository) getIOVariableByID(ctx context.Context, variableID string) (*Entities.IOVariable, error) {
	normalizedID := strings.TrimSpace(variableID)
	if normalizedID == "" {
		return nil, nil
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	res, err := r.adapter.Read(ioVariableTableName, map[string]string{"ID": normalizedID})
	if err != nil {
		return nil, err
	}

	record, err := firstRecord(res)
	if err != nil {
		return nil, nil
	}

	return recordToIOVariable(record)
}

func (r *TestCaseRepository) upsertIOVariable(ctx context.Context, variable Entities.IOVariable) error {
	variableID := strings.TrimSpace(variable.ID)
	if variableID == "" {
		return fmt.Errorf("io variable id is required")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return err
	}

	res, err := r.adapter.Read(ioVariableTableName, map[string]string{"ID": variableID})
	if err != nil {
		return err
	}

	if _, findErr := firstRecord(res); findErr != nil {
		_, insertErr := r.adapter.Insert(ioVariableTableName, []map[string]any{ioVariableToRecord(variable)})
		return insertErr
	}

	_, updateErr := r.adapter.Update(ioVariableTableName, "ID", variableID, ioVariableToUpdates(variable))
	return updateErr
}
