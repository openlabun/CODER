package challenge_entities

import consts "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/exam"

type IOVariable struct {
	ID 	   string					`json:"id"`
	Name   string					`json:"name"`
	Type   consts.VariableFormat	`json:"type"`
	Value  string					`json:"value"`
}