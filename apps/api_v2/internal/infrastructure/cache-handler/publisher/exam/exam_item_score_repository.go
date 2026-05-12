package cache_publisher

import (
	"context"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	ExamRepositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	cache_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
	PublisherAdapter "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/cache-handler/publisher"
)

type ExamItemScoreRepository struct {
	publisherAdapter *PublisherAdapter.CachePublisherAdapter
	repository       ExamRepositories.ExamItemScoreRepository
}

func NewExamItemScoreRepository(publisherAdapter *PublisherAdapter.CachePublisherAdapter, repository ExamRepositories.ExamItemScoreRepository) *ExamItemScoreRepository {
	return &ExamItemScoreRepository{publisherAdapter: publisherAdapter, repository: repository}
}

func (r *ExamItemScoreRepository) CreateExamItemScore(ctx context.Context, examItemScore *Entities.ExamItemScore) (*Entities.ExamItemScore, error) {
	message, err := PublisherAdapter.MapExamItemScoreToSyncMessage(examItemScore, cache_infrastructure.InsertOperation)
	if err != nil {
		return nil, err
	}
	if err := r.publisherAdapter.Save(message); err != nil {
		return nil, err
	}
	return examItemScore, nil
}

func (r *ExamItemScoreRepository) UpdateExamItemScore(ctx context.Context, examItemScore *Entities.ExamItemScore) (*Entities.ExamItemScore, error) {
	message, err := PublisherAdapter.MapExamItemScoreToSyncMessage(examItemScore, cache_infrastructure.UpdateOperation)
	if err != nil {
		return nil, err
	}
	if err := r.publisherAdapter.Save(message); err != nil {
		return nil, err
	}
	return examItemScore, nil
}

func (r *ExamItemScoreRepository) DeleteExamItemScore(ctx context.Context, examItemScoreID string) error {
	return r.repository.DeleteExamItemScore(ctx, examItemScoreID)
}

func (r *ExamItemScoreRepository) GetExamItemScoreByID(ctx context.Context, examItemScoreID string) (*Entities.ExamItemScore, error) {
	return r.repository.GetExamItemScoreByID(ctx, examItemScoreID)
}

func (r *ExamItemScoreRepository) GetExamItemScoresByExamScoreID(ctx context.Context, examScoreID string) ([]*Entities.ExamItemScore, error) {
	return r.repository.GetExamItemScoresByExamScoreID(ctx, examScoreID)
}

func (r *ExamItemScoreRepository) GetExamItemScore(ctx context.Context, examItemID, examScoreID string) (*Entities.ExamItemScore, error) {
	return r.repository.GetExamItemScore(ctx, examItemID, examScoreID)
}
