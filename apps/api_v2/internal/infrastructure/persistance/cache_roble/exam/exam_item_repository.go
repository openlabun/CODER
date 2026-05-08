package cache_roble_infrastructure

import (
	"context"
	"strings"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	cache "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
	roble_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
	repository "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/exam"
)

const cacheExamItemTable = "cache_exam_item"

type ExamItemRepository struct {
	robleRepository *repository.ExamItemRepository
	cacheAdapter    *cache.CacheAdapter
}

func NewExamItemRepository(robleAdapter *roble_infrastructure.RobleDatabaseAdapter, cacheAdapter *cache.CacheAdapter) *ExamItemRepository {
	return &ExamItemRepository{
		robleRepository: repository.NewExamItemRepository(robleAdapter),
		cacheAdapter:    cacheAdapter,
	}
}

func (r *ExamItemRepository) CreateExamItem(ctx context.Context, item *Entities.ExamItem) (*Entities.ExamItem, error) {
	result, err := r.robleRepository.CreateExamItem(ctx, item)
	if err != nil {
		return nil, err
	}
	if rec, e := cache.EntityToMap(result); e == nil {
		r.cacheAdapter.Save(cacheExamItemTable, "insert", rec)
	}
	return result, nil
}

func (r *ExamItemRepository) UpdateExamItem(ctx context.Context, item *Entities.ExamItem) (*Entities.ExamItem, error) {
	result, err := r.robleRepository.UpdateExamItem(ctx, item)
	if err != nil {
		return nil, err
	}
	if rec, e := cache.EntityToMap(result); e == nil {
		r.cacheAdapter.Save(cacheExamItemTable, "update", rec)
	}
	return result, nil
}

func (r *ExamItemRepository) DeleteExamItem(ctx context.Context, examItemID string) error {
	if err := r.robleRepository.DeleteExamItem(ctx, examItemID); err != nil {
		return err
	}
	r.cacheAdapter.DeleteByID(cacheExamItemTable, strings.TrimSpace(examItemID))
	return nil
}

func (r *ExamItemRepository) GetExamItemByID(ctx context.Context, examItemID string) (*Entities.ExamItem, error) {
	id := strings.TrimSpace(examItemID)
	if rec, err := r.cacheAdapter.FindByID(cacheExamItemTable, id); err == nil && rec != nil {
		if item, e := cache.MapToEntity[Entities.ExamItem](rec); e == nil {
			return item, nil
		}
	}
	item, err := r.robleRepository.GetExamItemByID(ctx, examItemID)
	if err != nil || item == nil {
		return item, err
	}
	if rec, e := cache.EntityToMap(item); e == nil {
		r.cacheAdapter.Save(cacheExamItemTable, "insert", rec)
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
		if recs, err := r.cacheAdapter.FindBy(cacheExamItemTable, conditions); err == nil && recs != nil {
			return examItemsFromRecords(recs), nil
		}
	}
	items, err := r.robleRepository.GetExamItem(ctx, examID, challengeID)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if rec, e := cache.EntityToMap(item); e == nil {
			r.cacheAdapter.Save(cacheExamItemTable, "insert", rec)
		}
	}
	return items, nil
}
