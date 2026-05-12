package cache_publisher

import (
	"context"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	ExamRepositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	cache_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
	PublisherAdapter "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/cache-handler/publisher"
)

type IOVariableRepository struct {
	publisherAdapter *PublisherAdapter.CachePublisherAdapter
	repository       ExamRepositories.IOVariableRepository
}

func NewIOVariableRepository(publisherAdapter *PublisherAdapter.CachePublisherAdapter, repository ExamRepositories.IOVariableRepository) *IOVariableRepository {
	return &IOVariableRepository{publisherAdapter: publisherAdapter, repository: repository}
}

func (r *IOVariableRepository) CreateIOVariable(ctx context.Context, ioVariable *Entities.IOVariable) (*Entities.IOVariable, error) {
	message, err := PublisherAdapter.MapIOVariableToSyncMessage(ioVariable, cache_infrastructure.InsertOperation)
	if err != nil {
		return nil, err
	}
	if err := r.publisherAdapter.Save(message); err != nil {
		return nil, err
	}
	return ioVariable, nil
}

func (r *IOVariableRepository) UpdateIOVariable(ctx context.Context, ioVariable *Entities.IOVariable) (*Entities.IOVariable, error) {
	message, err := PublisherAdapter.MapIOVariableToSyncMessage(ioVariable, cache_infrastructure.UpdateOperation)
	if err != nil {
		return nil, err
	}
	if err := r.publisherAdapter.Save(message); err != nil {
		return nil, err
	}
	return ioVariable, nil
}

func (r *IOVariableRepository) GetIOVariableByID(ctx context.Context, ioVariableID string) (*Entities.IOVariable, error) {
	return r.repository.GetIOVariableByID(ctx, ioVariableID)
}

func (r *IOVariableRepository) DeleteIOVariable(ctx context.Context, ioVariableID string) error {
	return r.repository.DeleteIOVariable(ctx, ioVariableID)
}
