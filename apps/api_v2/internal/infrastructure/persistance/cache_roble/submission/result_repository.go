package cache_roble_infrastructure

import (
	"context"
	"strings"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	cache "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
	roble_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/exam"
	repository "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/submission"
)

const cacheSubmissionResultTable = "cache_submission_result"

type SubmissionResultRepository struct {
	robleRepository *repository.SubmissionResultRepository
	cacheAdapter    *cache.CacheAdapter
}

func NewSubmissionResultRepository(robleAdapter *roble_infrastructure.RobleDatabaseAdapter, cacheAdapter *cache.CacheAdapter) *SubmissionResultRepository {
	ioVariableRepository := examRepository.NewIOVariableRepository(robleAdapter)
	return &SubmissionResultRepository{
		robleRepository: repository.NewSubmissionResultRepository(robleAdapter, ioVariableRepository),
		cacheAdapter:    cacheAdapter,
	}
}

func (r *SubmissionResultRepository) CreateResult(ctx context.Context, result *Entities.SubmissionResult) (*Entities.SubmissionResult, error) {
	res, err := r.robleRepository.CreateResult(ctx, result)
	if err != nil {
		return nil, err
	}
	if rec, e := cache.EntityToMap(res); e == nil {
		r.cacheAdapter.Save(cacheSubmissionResultTable, "insert", rec)
	}
	return res, nil
}

func (r *SubmissionResultRepository) UpdateResult(ctx context.Context, result *Entities.SubmissionResult) (*Entities.SubmissionResult, error) {
	res, err := r.robleRepository.UpdateResult(ctx, result)
	if err != nil {
		return nil, err
	}
	if rec, e := cache.EntityToMap(res); e == nil {
		r.cacheAdapter.Save(cacheSubmissionResultTable, "update", rec)
	}
	return res, nil
}

func (r *SubmissionResultRepository) DeleteResult(ctx context.Context, resultID string) error {
	if err := r.robleRepository.DeleteResult(ctx, resultID); err != nil {
		return err
	}
	r.cacheAdapter.DeleteByID(cacheSubmissionResultTable, strings.TrimSpace(resultID))
	return nil
}

func (r *SubmissionResultRepository) GetResultByID(ctx context.Context, resultID string) (*Entities.SubmissionResult, error) {
	id := strings.TrimSpace(resultID)
	if rec, err := r.cacheAdapter.FindByID(cacheSubmissionResultTable, id); err == nil && rec != nil {
		if res, e := cache.MapToEntity[Entities.SubmissionResult](rec); e == nil {
			return res, nil
		}
	}
	res, err := r.robleRepository.GetResultByID(ctx, resultID)
	if err != nil || res == nil {
		return res, err
	}
	if rec, e := cache.EntityToMap(res); e == nil {
		r.cacheAdapter.Save(cacheSubmissionResultTable, "insert", rec)
	}
	return res, nil
}

func (r *SubmissionResultRepository) GetResultsBySubmissionID(ctx context.Context, submissionID string) ([]*Entities.SubmissionResult, error) {
	id := strings.TrimSpace(submissionID)
	if recs, err := r.cacheAdapter.FindBy(cacheSubmissionResultTable, map[string]string{"submission_id": id}); err == nil && recs != nil {
		return resultsFromRecords(recs), nil
	}
	results, err := r.robleRepository.GetResultsBySubmissionID(ctx, submissionID)
	if err != nil {
		return nil, err
	}
	for _, res := range results {
		if rec, e := cache.EntityToMap(res); e == nil {
			r.cacheAdapter.Save(cacheSubmissionResultTable, "insert", rec)
		}
	}
	return results, nil
}

func (r *SubmissionResultRepository) GetResultByTestCase(ctx context.Context, testCaseID string) ([]*Entities.SubmissionResult, error) {
	id := strings.TrimSpace(testCaseID)
	if recs, err := r.cacheAdapter.FindBy(cacheSubmissionResultTable, map[string]string{"test_case_id": id}); err == nil && recs != nil {
		return resultsFromRecords(recs), nil
	}
	results, err := r.robleRepository.GetResultByTestCase(ctx, testCaseID)
	if err != nil {
		return nil, err
	}
	for _, res := range results {
		if rec, e := cache.EntityToMap(res); e == nil {
			r.cacheAdapter.Save(cacheSubmissionResultTable, "insert", rec)
		}
	}
	return results, nil
}