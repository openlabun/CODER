package cache_publisher

import (
	"context"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	SubmissionRepositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
	cache_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
	PublisherAdapter "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/cache-handler/publisher"
)

type SubmissionResultRepository struct {
	publisherAdapter *PublisherAdapter.CachePublisherAdapter
	repository       SubmissionRepositories.SubmissionResultRepository
}

func NewSubmissionResultRepository(publisherAdapter *PublisherAdapter.CachePublisherAdapter, repository SubmissionRepositories.SubmissionResultRepository) *SubmissionResultRepository {
	return &SubmissionResultRepository{publisherAdapter: publisherAdapter, repository: repository}
}

func (r *SubmissionResultRepository) CreateResult(ctx context.Context, result *Entities.SubmissionResult) (*Entities.SubmissionResult, error) {
	message, err := PublisherAdapter.MapSubmissionResultToSyncMessage(result, cache_infrastructure.InsertOperation)
	if err != nil {
		return nil, err
	}
	if err := r.publisherAdapter.Save(message); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *SubmissionResultRepository) UpdateResult(ctx context.Context, result *Entities.SubmissionResult) (*Entities.SubmissionResult, error) {
	message, err := PublisherAdapter.MapSubmissionResultToSyncMessage(result, cache_infrastructure.UpdateOperation)
	if err != nil {
		return nil, err
	}
	if err := r.publisherAdapter.Save(message); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *SubmissionResultRepository) DeleteResult(ctx context.Context, resultID string) error {
	return r.repository.DeleteResult(ctx, resultID)
}

func (r *SubmissionResultRepository) GetResultByID(ctx context.Context, resultID string) (*Entities.SubmissionResult, error) {
	return r.repository.GetResultByID(ctx, resultID)
}

func (r *SubmissionResultRepository) GetResultsBySubmissionID(ctx context.Context, submissionID string) ([]*Entities.SubmissionResult, error) {
	return r.repository.GetResultsBySubmissionID(ctx, submissionID)
}

func (r *SubmissionResultRepository) GetResultByTestCase(ctx context.Context, testCaseID string) ([]*Entities.SubmissionResult, error) {
	return r.repository.GetResultByTestCase(ctx, testCaseID)
}
