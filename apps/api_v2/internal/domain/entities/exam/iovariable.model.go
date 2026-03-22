package challenge_entities

type VariableFormat string
const (
	VariableFormatString VariableFormat = "string"
	VariableFormatInt    VariableFormat = "int"
	VariableFormatFloat  VariableFormat = "float"
)

type IOVariable struct {
	ID 	   string
	Name   string
	Type   VariableFormat
	Value  string
}