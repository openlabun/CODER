package cache_publisher

import (
	"context"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	ExamRepositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	cache_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
	PublisherAdapter "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/cache-handler/publisher"
)

type ExamItemRepository struct {
	publisherAdapter *PublisherAdapter.CachePublisherAdapter
	repository       ExamRepositories.ExamItemRepository
}

func NewExamItemRepository(publisherAdapter *PublisherAdapter.CachePublisherAdapter, repository ExamRepositories.ExamItemRepository) *ExamItemRepository {
	return &ExamItemRepository{publisherAdapter: publisherAdapter, repository: repository}
}

func (r *ExamItemRepository) CreateExamItem(ctx context.Context, examItem *Entities.ExamItem) (*Entities.ExamItem, error) {
	message, err := PublisherAdapter.MapExamItemToSyncMessage(examItem, cache_infrastructure.InsertOperation)
	if err != nil {
		return nil, err
	}
	if err := r.publisherAdapter.Save(message); err != nil {
		return nil, err
	}
	return examItem, nil
}

func (r *ExamItemRepository) UpdateExamItem(ctx context.Context, examItem *Entities.ExamItem) (*Entities.ExamItem, error) {
	message, err := PublisherAdapter.MapExamItemToSyncMessage(examItem, cache_infrastructure.UpdateOperation)
	if err != nil {
		return nil, err
	}
	if err := r.publisherAdapter.Save(message); err != nil {
		return nil, err
	}
	return examItem, nil
}

func (r *ExamItemRepository) DeleteExamItem(ctx context.Context, examItemID string) error {
	return r.repository.DeleteExamItem(ctx, examItemID)
}

func (r *ExamItemRepository) GetExamItemByID(ctx context.Context, examItemID string) (*Entities.ExamItem, error) {
	return r.repository.GetExamItemByID(ctx, examItemID)
}

func (r *ExamItemRepository) GetExamItem(ctx context.Context, examID *string, challengeID *string) ([]*Entities.ExamItem, error) {
	return r.repository.GetExamItem(ctx, examID, challengeID)
}
