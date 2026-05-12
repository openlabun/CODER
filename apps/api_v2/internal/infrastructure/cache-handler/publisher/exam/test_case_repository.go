package cache_publisher

import (
	"context"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	ExamRepositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	cache_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
	PublisherAdapter "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/cache-handler/publisher"
)

type TestCaseRepository struct {
	publisherAdapter *PublisherAdapter.CachePublisherAdapter
	repository       ExamRepositories.TestCaseRepository
}

func NewTestCaseRepository(publisherAdapter *PublisherAdapter.CachePublisherAdapter, repository ExamRepositories.TestCaseRepository) *TestCaseRepository {
	return &TestCaseRepository{publisherAdapter: publisherAdapter, repository: repository}
}

func (r *TestCaseRepository) CreateTestCase(ctx context.Context, testCase *Entities.TestCase) (*Entities.TestCase, error) {
	message, err := PublisherAdapter.MapTestCaseToSyncMessage(testCase, cache_infrastructure.InsertOperation)
	if err != nil {
		return nil, err
	}
	if err := r.publisherAdapter.Save(message); err != nil {
		return nil, err
	}
	return testCase, nil
}

func (r *TestCaseRepository) UpdateTestCase(ctx context.Context, testCase *Entities.TestCase) (*Entities.TestCase, error) {
	message, err := PublisherAdapter.MapTestCaseToSyncMessage(testCase, cache_infrastructure.UpdateOperation)
	if err != nil {
		return nil, err
	}
	if err := r.publisherAdapter.Save(message); err != nil {
		return nil, err
	}
	return testCase, nil
}

func (r *TestCaseRepository) DeleteTestCase(ctx context.Context, testCaseID string) error {
	return r.repository.DeleteTestCase(ctx, testCaseID)
}

func (r *TestCaseRepository) GetTestCaseByID(ctx context.Context, testCaseID string) (*Entities.TestCase, error) {
	return r.repository.GetTestCaseByID(ctx, testCaseID)
}

func (r *TestCaseRepository) GetTestCasesByChallengeID(ctx context.Context, challengeID string) ([]*Entities.TestCase, error) {
	return r.repository.GetTestCasesByChallengeID(ctx, challengeID)
}

func (r *TestCaseRepository) GetInputVariablesByTestCaseID(ctx context.Context, testCaseID string) ([]*Entities.IOVariable, error) {
	return r.repository.GetInputVariablesByTestCaseID(ctx, testCaseID)
}

func (r *TestCaseRepository) GetOutputVariablesByTestCaseID(ctx context.Context, testCaseID string) ([]*Entities.IOVariable, error) {
	return r.repository.GetOutputVariablesByTestCaseID(ctx, testCaseID)
}
