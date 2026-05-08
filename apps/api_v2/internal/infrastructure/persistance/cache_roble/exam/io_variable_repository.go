package cache_roble_infrastructure

import (
	"context"
	"strings"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	cache "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
	roble_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
	repository "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/exam"
)

const cacheIOVariableTable = "cache_io_variable"

type IOVariableRepository struct {
	robleRepository *repository.IOVariableRepository
	cacheAdapter    *cache.CacheAdapter
}

func NewIOVariableRepository(robleAdapter *roble_infrastructure.RobleDatabaseAdapter, cacheAdapter *cache.CacheAdapter) *IOVariableRepository {
	return &IOVariableRepository{
		robleRepository: repository.NewIOVariableRepository(robleAdapter),
		cacheAdapter:    cacheAdapter,
	}
}

func (r *IOVariableRepository) GetIOVariableByID(ctx context.Context, variableID string) (*Entities.IOVariable, error) {
	id := strings.TrimSpace(variableID)
	if rec, err := r.cacheAdapter.FindByID(cacheIOVariableTable, id); err == nil && rec != nil {
		if v, e := cache.MapToEntity[Entities.IOVariable](rec); e == nil {
			return v, nil
		}
	}
	v, err := r.robleRepository.GetIOVariableByID(ctx, variableID)
	if err != nil || v == nil {
		return v, err
	}
	if rec, e := cache.EntityToMap(v); e == nil {
		r.cacheAdapter.Save(cacheIOVariableTable, "insert", rec)
	}
	return v, nil
}

func (r *IOVariableRepository) CreateIOVariable(ctx context.Context, variable *Entities.IOVariable) (*Entities.IOVariable, error) {
	result, err := r.robleRepository.CreateIOVariable(ctx, variable)
	if err != nil {
		return nil, err
	}
	if rec, e := cache.EntityToMap(result); e == nil {
		r.cacheAdapter.Save(cacheIOVariableTable, "insert", rec)
	}
	return result, nil
}

func (r *IOVariableRepository) DeleteIOVariable(ctx context.Context, variableID string) error {
	if err := r.robleRepository.DeleteIOVariable(ctx, variableID); err != nil {
		return err
	}
	r.cacheAdapter.DeleteByID(cacheIOVariableTable, strings.TrimSpace(variableID))
	return nil
}
