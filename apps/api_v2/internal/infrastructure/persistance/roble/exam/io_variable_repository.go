package roble_infrastructure

import (
	"context"
	"fmt"
	"strings"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
)

const (
	ioVariableTableName = "IOVariable"
)

type IOVariableRepository struct {
	adapter *infrastructure.RobleDatabaseAdapter
}

func NewIOVariableRepository(adapter *infrastructure.RobleDatabaseAdapter) *IOVariableRepository {
	return &IOVariableRepository{adapter: adapter}
}

func (r *IOVariableRepository) CreateIOVariable(ctx context.Context, ioVariable *Entities.IOVariable) (*Entities.IOVariable, error) {
	if ioVariable == nil {
		return nil, fmt.Errorf("io variable is nil")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	ioVariableID := strings.TrimSpace(ioVariable.ID)
	if ioVariableID == "" {
		return nil, fmt.Errorf("io variable id is required")
	}

	_, err := r.adapter.Insert(ioVariableTableName, []map[string]any{ioVariableToRecord(*ioVariable)})
	if err != nil {
		return nil, err
	}

	return r.GetIOVariableByID(ctx, ioVariableID)
}

func (r *IOVariableRepository) UpdateIOVariable(ctx context.Context, ioVariable *Entities.IOVariable) (*Entities.IOVariable, error) {
	if ioVariable == nil {
		return nil, fmt.Errorf("io variable is nil")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	ioVariableID := strings.TrimSpace(ioVariable.ID)
	if ioVariableID == "" {
		return nil, fmt.Errorf("io variable id is required")
	}

	_, err := r.adapter.Update(ioVariableTableName, "ID", ioVariableID, ioVariableToUpdates(*ioVariable))
	if err != nil {
		return nil, err
	}

	return r.GetIOVariableByID(ctx, ioVariableID)
}

func (r *IOVariableRepository) GetIOVariableByID(ctx context.Context, ioVariableID string) (*Entities.IOVariable, error) {
	normalizedID := strings.TrimSpace(ioVariableID)
	if normalizedID == "" {
		return nil, fmt.Errorf("ioVariableID is required")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	res, err := r.adapter.Read(ioVariableTableName, map[string]string{"ID": normalizedID})
	if err != nil {
		return nil, err
	}

	record, err := firstRecord(res)
	if err != nil {
		return nil, nil
	}

	return recordToIOVariable(record)
}

func (r *IOVariableRepository) DeleteIOVariable(ctx context.Context, ioVariableID string) error {
	normalizedID := strings.TrimSpace(ioVariableID)
	if normalizedID == "" {
		return fmt.Errorf("ioVariableID is required")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return err
	}

	_, err := r.adapter.Delete(ioVariableTableName, "ID", normalizedID)
	return err
}
