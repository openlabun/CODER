package cache_roble_infrastructure

import (
	"context"
	"strings"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	Repository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	cache "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
)

type TestCaseRepository struct {
	repository   Repository.TestCaseRepository
	cacheAdapter *cache.CacheAdapter
}

func NewTestCaseRepository(repository Repository.TestCaseRepository, cacheAdapter *cache.CacheAdapter) *TestCaseRepository {
	return &TestCaseRepository{
		repository:   repository,
		cacheAdapter: cacheAdapter,
	}
}

func (r *TestCaseRepository) CreateTestCase(ctx context.Context, testCase *Entities.TestCase) (*Entities.TestCase, error) {
	result, err := r.repository.CreateTestCase(ctx, testCase)
	if err != nil {
		return nil, err
	}
	if rec, e := cache.EntityToMap(result); e == nil {
		r.cacheAdapter.Save(cache.CacheTestCaseTable, cache.InsertOperation, rec)
	}
	return result, nil
}

func (r *TestCaseRepository) UpdateTestCase(ctx context.Context, testCase *Entities.TestCase) (*Entities.TestCase, error) {
	result, err := r.repository.UpdateTestCase(ctx, testCase)
	if err != nil {
		return nil, err
	}
	if rec, e := cache.EntityToMap(result); e == nil {
		r.cacheAdapter.Save(cache.CacheTestCaseTable, cache.UpdateOperation, rec)
	}
	return result, nil
}

func (r *TestCaseRepository) DeleteTestCase(ctx context.Context, testCaseID string) error {
	if err := r.repository.DeleteTestCase(ctx, testCaseID); err != nil {
		return err
	}
	r.cacheAdapter.DeleteByID(cache.CacheTestCaseTable, strings.TrimSpace(testCaseID))
	return nil
}

func (r *TestCaseRepository) GetTestCaseByID(ctx context.Context, testCaseID string) (*Entities.TestCase, error) {
	id := strings.TrimSpace(testCaseID)
	if rec, err := r.cacheAdapter.FindByID(cache.CacheTestCaseTable, id); err == nil && rec != nil {
		cache.NormalizeBools(rec, "is_sample", "custom")
		if t, e := cache.MapToEntity[Entities.TestCase](rec); e == nil {
			return t, nil
		}
	}
	t, err := r.repository.GetTestCaseByID(ctx, testCaseID)
	if err != nil || t == nil {
		return t, err
	}
	if rec, e := cache.EntityToMap(t); e == nil {
		r.cacheAdapter.Save(cache.CacheTestCaseTable, cache.InsertOperation, rec)
	}
	return t, nil
}

func (r *TestCaseRepository) GetTestCasesByChallengeID(ctx context.Context, challengeID string) ([]*Entities.TestCase, error) {
	id := strings.TrimSpace(challengeID)
	if recs, err := r.cacheAdapter.FindBy(cache.CacheTestCaseTable, map[string]string{"challenge_id": id}); err == nil && recs != nil {
		return testCasesFromRecords(recs), nil
	}
	testCases, err := r.repository.GetTestCasesByChallengeID(ctx, challengeID)
	if err != nil {
		return nil, err
	}
	for _, t := range testCases {
		if rec, e := cache.EntityToMap(t); e == nil {
			r.cacheAdapter.Save(cache.CacheTestCaseTable, cache.InsertOperation, rec)
		}
	}
	return testCases, nil
}

func (r *TestCaseRepository) GetInputVariablesByTestCaseID(ctx context.Context, testCaseID string) ([]*Entities.IOVariable, error) {
	t, err := r.GetTestCaseByID(ctx, testCaseID)
	if err != nil || t == nil {
		return []*Entities.IOVariable{}, err
	}
	vars := make([]*Entities.IOVariable, 0, len(t.Input))
	for i := range t.Input {
		v := t.Input[i]
		vars = append(vars, &v)
	}
	return vars, nil
}

func (r *TestCaseRepository) GetOutputVariablesByTestCaseID(ctx context.Context, testCaseID string) ([]*Entities.IOVariable, error) {
	t, err := r.GetTestCaseByID(ctx, testCaseID)
	if err != nil || t == nil || strings.TrimSpace(t.ExpectedOutput.ID) == "" {
		return []*Entities.IOVariable{}, err
	}
	v := t.ExpectedOutput
	return []*Entities.IOVariable{&v}, nil
}