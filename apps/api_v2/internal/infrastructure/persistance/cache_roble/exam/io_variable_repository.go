package cache_roble_infrastructure

import (
	"context"
	"strings"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	Repository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	cache "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
)

type IOVariableRepository struct {
	repository   Repository.IOVariableRepository
	cacheAdapter *cache.CacheAdapter
}

func NewIOVariableRepository(repository Repository.IOVariableRepository, cacheAdapter *cache.CacheAdapter) *IOVariableRepository {
	return &IOVariableRepository{
		repository:   repository,
		cacheAdapter: cacheAdapter,
	}
}

func (r *IOVariableRepository) GetIOVariableByID(ctx context.Context, variableID string) (*Entities.IOVariable, error) {
	id := strings.TrimSpace(variableID)
	if rec, err := r.cacheAdapter.FindByID(cache.CacheIOVariableTable, id); err == nil && rec != nil {
		if v, e := cache.MapToEntity[Entities.IOVariable](rec); e == nil {
			return v, nil
		}
	}
	v, err := r.repository.GetIOVariableByID(ctx, variableID)
	if err != nil || v == nil {
		return v, err
	}
	if rec, e := cache.EntityToMap(v); e == nil {
		r.cacheAdapter.Save(cache.CacheIOVariableTable, cache.InsertOperation, rec)
	}
	return v, nil
}

func (r *IOVariableRepository) CreateIOVariable(ctx context.Context, variable *Entities.IOVariable) (*Entities.IOVariable, error) {
	result, err := r.repository.CreateIOVariable(ctx, variable)
	if err != nil {
		return nil, err
	}
	if rec, e := cache.EntityToMap(result); e == nil {
		r.cacheAdapter.Save(cache.CacheIOVariableTable, cache.InsertOperation, rec)
	}
	return result, nil
}

func (r *IOVariableRepository) UpdateIOVariable(ctx context.Context, variable *Entities.IOVariable) (*Entities.IOVariable, error) {
	result, err := r.repository.UpdateIOVariable(ctx, variable)
	if err != nil {
		return nil, err
	}
	if rec, e := cache.EntityToMap(result); e == nil {
		r.cacheAdapter.Save(cache.CacheIOVariableTable, cache.UpdateOperation, rec)
	}
	return result, nil
}

func (r *IOVariableRepository) DeleteIOVariable(ctx context.Context, variableID string) error {
	if err := r.repository.DeleteIOVariable(ctx, variableID); err != nil {
		return err
	}
	r.cacheAdapter.DeleteByID(cache.CacheIOVariableTable, strings.TrimSpace(variableID))
	return nil
}
