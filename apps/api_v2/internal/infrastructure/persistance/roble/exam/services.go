package roble_infrastructure

import (
	"context"
	"strings"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
)


func deleteIOVariablesByIDs(ctx context.Context, adapter *infrastructure.RobleDatabaseAdapter, ids []string) error {
	if adapter == nil || len(ids) == 0 {
		return nil
	}

	seen := map[string]struct{}{}
	for _, id := range ids {
		normalized := strings.TrimSpace(id)
		if normalized == "" {
			continue
		}
		if _, exists := seen[normalized]; exists {
			continue
		}
		seen[normalized] = struct{}{}
		if err := infrastructure.SetAdapterTokenFromContext(ctx, adapter); err != nil {
			return err
		}

		if _, err := adapter.Delete(ioVariableTableName, "ID", normalized); err != nil {
			return err
		}
	}

	return nil
}

func ioVariableIDs(variables []Entities.IOVariable) []string {
	ids := make([]string, 0, len(variables))
	for _, variable := range variables {
		if id := strings.TrimSpace(variable.ID); id != "" {
			ids = append(ids, id)
		}
	}

	return ids
}

func relatedIOVariableIDs(record map[string]any) []string {
	if record == nil {
		return nil
	}

	seen := map[string]struct{}{}
	ids := make([]string, 0)

	for _, id := range asStringList(record["Input"]) {
		normalized := strings.TrimSpace(id)
		if normalized == "" {
			continue
		}
		if _, exists := seen[normalized]; exists {
			continue
		}
		seen[normalized] = struct{}{}
		ids = append(ids, normalized)
	}

	for _, id := range asStringList(record["InputVariables"]) {
		normalized := strings.TrimSpace(id)
		if normalized == "" {
			continue
		}
		if _, exists := seen[normalized]; exists {
			continue
		}
		seen[normalized] = struct{}{}
		ids = append(ids, normalized)
	}

	if outputID := strings.TrimSpace(asString(record["Output"])); outputID != "" {
		if _, exists := seen[outputID]; !exists {
			seen[outputID] = struct{}{}
			ids = append(ids, outputID)
		}
	}

	if outputID := strings.TrimSpace(asString(record["ExpectedOutput"])); outputID != "" {
		if _, exists := seen[outputID]; !exists {
			seen[outputID] = struct{}{}
			ids = append(ids, outputID)
		}
	}

	if outputID := strings.TrimSpace(asString(record["OutputVariable"])); outputID != "" {
		if _, exists := seen[outputID]; !exists {
			seen[outputID] = struct{}{}
			ids = append(ids, outputID)
		}
	}

	return ids
}
