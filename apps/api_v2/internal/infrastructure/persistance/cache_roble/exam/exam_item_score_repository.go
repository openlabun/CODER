package cache_roble_infrastructure

import (
	"context"
	"strings"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	cache "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
	roble_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
	repository "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/exam"
)

type ExamItemScoreRepository struct {
	robleRepository *repository.ExamItemScoreRepository
	cacheAdapter    *cache.CacheAdapter
}

func NewExamItemScoreRepository(robleAdapter *roble_infrastructure.RobleDatabaseAdapter, cacheAdapter *cache.CacheAdapter) *ExamItemScoreRepository {
	return &ExamItemScoreRepository{
		robleRepository: repository.NewExamItemScoreRepository(robleAdapter),
		cacheAdapter:    cacheAdapter,
	}
}

func (r *ExamItemScoreRepository) CreateExamItemScore(ctx context.Context, s *Entities.ExamItemScore) (*Entities.ExamItemScore, error) {
	result, err := r.robleRepository.CreateExamItemScore(ctx, s)
	if err != nil {
		return nil, err
	}
	if rec, e := cache.EntityToMap(result); e == nil {
		r.cacheAdapter.Save(cache.CacheExamItemScoreTable, cache.InsertOperation, rec)
	}
	return result, nil
}

func (r *ExamItemScoreRepository) UpdateExamItemScore(ctx context.Context, s *Entities.ExamItemScore) (*Entities.ExamItemScore, error) {
	result, err := r.robleRepository.UpdateExamItemScore(ctx, s)
	if err != nil {
		return nil, err
	}
	if rec, e := cache.EntityToMap(result); e == nil {
		r.cacheAdapter.Save(cache.CacheExamItemScoreTable, cache.UpdateOperation, rec)
	}
	return result, nil
}

func (r *ExamItemScoreRepository) DeleteExamItemScore(ctx context.Context, examItemScoreID string) error {
	if err := r.robleRepository.DeleteExamItemScore(ctx, examItemScoreID); err != nil {
		return err
	}
	r.cacheAdapter.DeleteByID(cache.CacheExamItemScoreTable, strings.TrimSpace(examItemScoreID))
	return nil
}

func (r *ExamItemScoreRepository) GetExamItemScoreByID(ctx context.Context, id string) (*Entities.ExamItemScore, error) {
	if rec, err := r.cacheAdapter.FindByID(cache.CacheExamItemScoreTable, strings.TrimSpace(id)); err == nil && rec != nil {
		if s, e := cache.MapToEntity[Entities.ExamItemScore](rec); e == nil {
			return s, nil
		}
	}
	s, err := r.robleRepository.GetExamItemScoreByID(ctx, id)
	if err != nil || s == nil {
		return s, err
	}
	if rec, e := cache.EntityToMap(s); e == nil {
		r.cacheAdapter.Save(cache.CacheExamItemScoreTable, cache.InsertOperation, rec)
	}
	return s, nil
}

func (r *ExamItemScoreRepository) GetExamItemScore(ctx context.Context, examItemID, examScoreID string) (*Entities.ExamItemScore, error) {
	conditions := map[string]string{
		"exam_item_id":  strings.TrimSpace(examItemID),
		"exam_score_id": strings.TrimSpace(examScoreID),
	}
	if recs, err := r.cacheAdapter.FindBy(cache.CacheExamItemScoreTable, conditions); err == nil && len(recs) > 0 {
		if s, e := cache.MapToEntity[Entities.ExamItemScore](recs[0]); e == nil {
			return s, nil
		}
	}
	s, err := r.robleRepository.GetExamItemScore(ctx, examItemID, examScoreID)
	if err != nil || s == nil {
		return s, err
	}
	if rec, e := cache.EntityToMap(s); e == nil {
		r.cacheAdapter.Save(cache.CacheExamItemScoreTable, cache.InsertOperation, rec)
	}
	return s, nil
}

func (r *ExamItemScoreRepository) GetExamItemScoresByExamScoreID(ctx context.Context, examScoreID string) ([]*Entities.ExamItemScore, error) {
	id := strings.TrimSpace(examScoreID)
	if recs, err := r.cacheAdapter.FindBy(cache.CacheExamItemScoreTable, map[string]string{"exam_score_id": id}); err == nil && recs != nil {
		return examItemScoresFromRecords(recs), nil
	}
	scores, err := r.robleRepository.GetExamItemScoresByExamScoreID(ctx, examScoreID)
	if err != nil {
		return nil, err
	}
	for _, s := range scores {
		if rec, e := cache.EntityToMap(s); e == nil {
			r.cacheAdapter.Save(cache.CacheExamItemScoreTable, cache.InsertOperation, rec)
		}
	}
	return scores, nil
}