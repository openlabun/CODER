package roble_infrastructure

import (
	"strings"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
)

func ioVariableIDs(variables []Entities.IOVariable) []string {
	ids := make([]string, 0, len(variables))
	for _, variable := range variables {
		if id := strings.TrimSpace(variable.ID); id != "" {
			ids = append(ids, id)
		}
	}

	return ids
}
