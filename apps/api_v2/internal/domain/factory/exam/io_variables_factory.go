package exam_factory

import (
	"strings"

	"github.com/google/uuid"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	Validations "github.com/openlabun/CODER/apps/api_v2/internal/domain/validations/exam"
)

func NewIOVariable(name string, variableType Entities.VariableFormat, value string) (*Entities.IOVariable, error) {
	ioVariable := &Entities.IOVariable{
		ID:    uuid.New().String(),
		Name:  strings.TrimSpace(name),
		Type:  variableType,
		Value: strings.TrimSpace(value),
	}

	if err := Validations.ValidateIOVariable(*ioVariable); err != nil {
		return nil, err
	}

	return ioVariable, nil
}

func ExistingIOVariable(id, name string, variableType Entities.VariableFormat, value string) (*Entities.IOVariable, error) {
	ioVariable := &Entities.IOVariable{
		ID:    strings.TrimSpace(id),
		Name:  strings.TrimSpace(name),
		Type:  variableType,
		Value: strings.TrimSpace(value),
	}

	if err := Validations.ValidateIOVariable(*ioVariable); err != nil {
		return nil, err
	}

	return ioVariable, nil
}

