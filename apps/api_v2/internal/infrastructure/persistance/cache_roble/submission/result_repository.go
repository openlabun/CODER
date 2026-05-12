package cache_roble_infrastructure

import (
	"context"
	"strings"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	Repository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
	cache "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
)

type SubmissionResultRepository struct {
	repository   Repository.SubmissionResultRepository
	cacheAdapter *cache.CacheAdapter
}

func NewSubmissionResultRepository(repository Repository.SubmissionResultRepository, cacheAdapter *cache.CacheAdapter) *SubmissionResultRepository {
	return &SubmissionResultRepository{
		repository:   repository,
		cacheAdapter: cacheAdapter,
	}
}

func (r *SubmissionResultRepository) CreateResult(ctx context.Context, result *Entities.SubmissionResult) (*Entities.SubmissionResult, error) {
	res, err := r.repository.CreateResult(ctx, result)
	if err != nil {
		return nil, err
	}
	if rec, e := cache.EntityToMap(res); e == nil {
		r.cacheAdapter.Save(cache.CacheSubmissionResultTable,  cache.InsertOperation, rec)
	}
	return res, nil
}

func (r *SubmissionResultRepository) UpdateResult(ctx context.Context, result *Entities.SubmissionResult) (*Entities.SubmissionResult, error) {
	res, err := r.repository.UpdateResult(ctx, result)
	if err != nil {
		return nil, err
	}
	if rec, e := cache.EntityToMap(res); e == nil {
		r.cacheAdapter.Save(cache.CacheSubmissionResultTable, cache.UpdateOperation, rec)
	}
	return res, nil
}

func (r *SubmissionResultRepository) DeleteResult(ctx context.Context, resultID string) error {
	if err := r.repository.DeleteResult(ctx, resultID); err != nil {
		return err
	}
	r.cacheAdapter.DeleteByID(cache.CacheSubmissionResultTable, strings.TrimSpace(resultID))
	return nil
}

func (r *SubmissionResultRepository) GetResultByID(ctx context.Context, resultID string) (*Entities.SubmissionResult, error) {
	id := strings.TrimSpace(resultID)
	if rec, err := r.cacheAdapter.FindByID(cache.CacheSubmissionResultTable, id); err == nil && rec != nil {
		if res, e := cache.MapToEntity[Entities.SubmissionResult](rec); e == nil {
			return res, nil
		}
	}
	res, err := r.repository.GetResultByID(ctx, resultID)
	if err != nil || res == nil {
		return res, err
	}
	if rec, e := cache.EntityToMap(res); e == nil {
		r.cacheAdapter.Save(cache.CacheSubmissionResultTable, cache.InsertOperation, rec)
	}
	return res, nil
}

func (r *SubmissionResultRepository) GetResultsBySubmissionID(ctx context.Context, submissionID string) ([]*Entities.SubmissionResult, error) {
	id := strings.TrimSpace(submissionID)
	if recs, err := r.cacheAdapter.FindBy(cache.CacheSubmissionResultTable, map[string]string{"submission_id": id}); err == nil && recs != nil {
		return resultsFromRecords(recs), nil
	}
	results, err := r.repository.GetResultsBySubmissionID(ctx, submissionID)
	if err != nil {
		return nil, err
	}
	for _, res := range results {
		if rec, e := cache.EntityToMap(res); e == nil {
			r.cacheAdapter.Save(cache.CacheSubmissionResultTable,  cache.InsertOperation, rec)
		}
	}
	return results, nil
}

func (r *SubmissionResultRepository) GetResultByTestCase(ctx context.Context, testCaseID string) ([]*Entities.SubmissionResult, error) {
	id := strings.TrimSpace(testCaseID)
	if recs, err := r.cacheAdapter.FindBy(cache.CacheSubmissionResultTable, map[string]string{"test_case_id": id}); err == nil && recs != nil {
		return resultsFromRecords(recs), nil
	}
	results, err := r.repository.GetResultByTestCase(ctx, testCaseID)
	if err != nil {
		return nil, err
	}
	for _, res := range results {
		if rec, e := cache.EntityToMap(res); e == nil {
			r.cacheAdapter.Save(cache.CacheSubmissionResultTable,  cache.InsertOperation, rec)
		}
	}
	return results, nil
}