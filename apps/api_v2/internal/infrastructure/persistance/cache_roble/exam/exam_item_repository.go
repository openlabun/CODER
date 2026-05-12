package cache_roble_infrastructure

import (
	"context"
	"strings"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	Repository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	cache "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
)

type ExamItemRepository struct {
	repository   Repository.ExamItemRepository
	cacheAdapter *cache.CacheAdapter
}

func NewExamItemRepository(repository Repository.ExamItemRepository, cacheAdapter *cache.CacheAdapter) *ExamItemRepository {
	return &ExamItemRepository{
		repository:   repository,
		cacheAdapter: cacheAdapter,
	}
}

func (r *ExamItemRepository) CreateExamItem(ctx context.Context, item *Entities.ExamItem) (*Entities.ExamItem, error) {
	result, err := r.repository.CreateExamItem(ctx, item)
	if err != nil {
		return nil, err
	}
	if rec, e := cache.EntityToMap(result); e == nil {
		r.cacheAdapter.Save(cache.CacheExamItemTable, cache.InsertOperation, rec)
	}
	return result, nil
}

func (r *ExamItemRepository) UpdateExamItem(ctx context.Context, item *Entities.ExamItem) (*Entities.ExamItem, error) {
	result, err := r.repository.UpdateExamItem(ctx, item)
	if err != nil {
		return nil, err
	}
	if rec, e := cache.EntityToMap(result); e == nil {
		r.cacheAdapter.Save(cache.CacheExamItemTable, cache.UpdateOperation, rec)
	}
	return result, nil
}

func (r *ExamItemRepository) DeleteExamItem(ctx context.Context, examItemID string) error {
	if err := r.repository.DeleteExamItem(ctx, examItemID); err != nil {
		return err
	}
	r.cacheAdapter.DeleteByID(cache.CacheExamItemTable, strings.TrimSpace(examItemID))
	return nil
}

func (r *ExamItemRepository) GetExamItemByID(ctx context.Context, examItemID string) (*Entities.ExamItem, error) {
	id := strings.TrimSpace(examItemID)
	if rec, err := r.cacheAdapter.FindByID(cache.CacheExamItemTable, id); err == nil && rec != nil {
		if item, e := cache.MapToEntity[Entities.ExamItem](rec); e == nil {
			return item, nil
		}
	}
	item, err := r.repository.GetExamItemByID(ctx, examItemID)
	if err != nil || item == nil {
		return item, err
	}
	if rec, e := cache.EntityToMap(item); e == nil {
		r.cacheAdapter.Save(cache.CacheExamItemTable, cache.InsertOperation, rec)
	}
	return item, nil
}

func (r *ExamItemRepository) GetExamItem(ctx context.Context, examID *string, challengeID *string) ([]*Entities.ExamItem, error) {
	conditions := map[string]string{}
	if examID != nil && strings.TrimSpace(*examID) != "" {
		conditions["exam_id"] = strings.TrimSpace(*examID)
	}
	if challengeID != nil && strings.TrimSpace(*challengeID) != "" {
		conditions["challenge_id"] = strings.TrimSpace(*challengeID)
	}
	if len(conditions) > 0 {
		if recs, err := r.cacheAdapter.FindBy(cache.CacheExamItemTable, conditions); err == nil && recs != nil {
			return examItemsFromRecords(recs), nil
		}
	}
	items, err := r.repository.GetExamItem(ctx, examID, challengeID)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if rec, e := cache.EntityToMap(item); e == nil {
			r.cacheAdapter.Save(cache.CacheExamItemTable, cache.InsertOperation, rec)
		}
	}
	return items, nil
}
