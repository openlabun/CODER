package challenge_entities

import exam_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/exam"

type VariableFormat = exam_constants.VariableFormat

const (
	VariableFormatString VariableFormat = exam_constants.VariableFormatString
	VariableFormatInt    VariableFormat = exam_constants.VariableFormatInt
	VariableFormatFloat  VariableFormat = exam_constants.VariableFormatFloat
)

type IOVariable struct {
	ID 	   string			`json:"id"`
	Name   string			`json:"name"`
	Type   VariableFormat	`json:"type"`
	Value  string			`json:"value"`
}