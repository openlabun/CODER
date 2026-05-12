package cache_roble_infrastructure

import (
	"context"
	"strings"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	Repository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	cache "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
)

type ExamScoreRepository struct {
	repository   Repository.ExamScoreRepository
	cacheAdapter *cache.CacheAdapter
}

func NewExamScoreRepository(repository Repository.ExamScoreRepository, cacheAdapter *cache.CacheAdapter) *ExamScoreRepository {
	return &ExamScoreRepository{
		repository:   repository,
		cacheAdapter: cacheAdapter,
	}
}

func (r *ExamScoreRepository) CreateExamScore(ctx context.Context, s *Entities.ExamScore) (*Entities.ExamScore, error) {
	result, err := r.repository.CreateExamScore(ctx, s)
	if err != nil {
		return nil, err
	}
	if rec, e := cache.EntityToMap(result); e == nil {
		r.cacheAdapter.Save(cache.CacheExamScoreTable, cache.InsertOperation, rec)
	}
	return result, nil
}

func (r *ExamScoreRepository) UpdateExamScore(ctx context.Context, s *Entities.ExamScore) (*Entities.ExamScore, error) {
	result, err := r.repository.UpdateExamScore(ctx, s)
	if err != nil {
		return nil, err
	}
	if rec, e := cache.EntityToMap(result); e == nil {
		r.cacheAdapter.Save(cache.CacheExamScoreTable, cache.UpdateOperation, rec)
	}
	return result, nil
}

func (r *ExamScoreRepository) DeleteExamScore(ctx context.Context, examScoreID string) error {
	if err := r.repository.DeleteExamScore(ctx, examScoreID); err != nil {
		return err
	}
	r.cacheAdapter.DeleteByID(cache.CacheExamScoreTable, strings.TrimSpace(examScoreID))
	return nil
}

func (r *ExamScoreRepository) GetExamScoreByID(ctx context.Context, examScoreID string) (*Entities.ExamScore, error) {
	id := strings.TrimSpace(examScoreID)
	if rec, err := r.cacheAdapter.FindByID(cache.CacheExamScoreTable, id); err == nil && rec != nil {
		if s, e := cache.MapToEntity[Entities.ExamScore](rec); e == nil {
			return s, nil
		}
	}
	s, err := r.repository.GetExamScoreByID(ctx, examScoreID)
	if err != nil || s == nil {
		return s, err
	}
	if rec, e := cache.EntityToMap(s); e == nil {
		r.cacheAdapter.Save(cache.CacheExamScoreTable, cache.InsertOperation, rec)
	}
	return s, nil
}

func (r *ExamScoreRepository) GetExamScoreBySessionID(ctx context.Context, sessionID string) (*Entities.ExamScore, error) {
	id := strings.TrimSpace(sessionID)
	if recs, err := r.cacheAdapter.FindBy(cache.CacheExamScoreTable, map[string]string{"session_id": id}); err == nil && len(recs) > 0 {
		if s, e := cache.MapToEntity[Entities.ExamScore](recs[0]); e == nil {
			return s, nil
		}
	}
	s, err := r.repository.GetExamScoreBySessionID(ctx, sessionID)
	if err != nil || s == nil {
		return s, err
	}
	if rec, e := cache.EntityToMap(s); e == nil {
		r.cacheAdapter.Save(cache.CacheExamScoreTable, cache.InsertOperation, rec)
	}
	return s, nil
}

func (r *ExamScoreRepository) GetExamScores(ctx context.Context, examID, studentID *string) ([]*Entities.ExamScore, error) {
	conditions := map[string]string{}
	if examID != nil && strings.TrimSpace(*examID) != "" {
		conditions["exam_id"] = strings.TrimSpace(*examID)
	}
	if studentID != nil && strings.TrimSpace(*studentID) != "" {
		conditions["student_id"] = strings.TrimSpace(*studentID)
	}
	if len(conditions) > 0 {
		if recs, err := r.cacheAdapter.FindBy(cache.CacheExamScoreTable, conditions); err == nil && recs != nil {
			return examScoresFromRecords(recs), nil
		}
	}
	scores, err := r.repository.GetExamScores(ctx, examID, studentID)
	if err != nil {
		return nil, err
	}
	for _, s := range scores {
		if rec, e := cache.EntityToMap(s); e == nil {
			r.cacheAdapter.Save(cache.CacheExamScoreTable, cache.InsertOperation, rec)
		}
	}
	return scores, nil
}