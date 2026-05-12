package cache_publisher

import (
	"context"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	ExamRepositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	cache_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
	PublisherAdapter "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/cache-handler/publisher"
)

type ExamScoreRepository struct {
	publisherAdapter *PublisherAdapter.CachePublisherAdapter
	repository       ExamRepositories.ExamScoreRepository
}

func NewExamScoreRepository(publisherAdapter *PublisherAdapter.CachePublisherAdapter, repository ExamRepositories.ExamScoreRepository) *ExamScoreRepository {
	return &ExamScoreRepository{publisherAdapter: publisherAdapter, repository: repository}
}

func (r *ExamScoreRepository) CreateExamScore(ctx context.Context, examScore *Entities.ExamScore) (*Entities.ExamScore, error) {
	message, err := PublisherAdapter.MapExamScoreToSyncMessage(examScore, cache_infrastructure.InsertOperation)
	if err != nil {
		return nil, err
	}
	if err := r.publisherAdapter.Save(message); err != nil {
		return nil, err
	}
	return examScore, nil
}

func (r *ExamScoreRepository) UpdateExamScore(ctx context.Context, examScore *Entities.ExamScore) (*Entities.ExamScore, error) {
	message, err := PublisherAdapter.MapExamScoreToSyncMessage(examScore, cache_infrastructure.UpdateOperation)
	if err != nil {
		return nil, err
	}
	if err := r.publisherAdapter.Save(message); err != nil {
		return nil, err
	}
	return examScore, nil
}

func (r *ExamScoreRepository) DeleteExamScore(ctx context.Context, examScoreID string) error {
	return r.repository.DeleteExamScore(ctx, examScoreID)
}

func (r *ExamScoreRepository) GetExamScores(ctx context.Context, examID, studentID *string) ([]*Entities.ExamScore, error) {
	return r.repository.GetExamScores(ctx, examID, studentID)
}

func (r *ExamScoreRepository) GetExamScoreByID(ctx context.Context, examScoreID string) (*Entities.ExamScore, error) {
	return r.repository.GetExamScoreByID(ctx, examScoreID)
}

func (r *ExamScoreRepository) GetExamScoreBySessionID(ctx context.Context, sessionID string) (*Entities.ExamScore, error) {
	return r.repository.GetExamScoreBySessionID(ctx, sessionID)
}
