package challenge_entities

type VariableFormat string
const (
	VariableFormatString VariableFormat = "string"
	VariableFormatInt    VariableFormat = "int"
	VariableFormatFloat  VariableFormat = "float"
)

type IOVariable struct {
	ID 	   string			`json:"id"`
	Name   string			`json:"name"`
	Type   VariableFormat	`json:"type"`
	Value  string			`json:"value"`
}