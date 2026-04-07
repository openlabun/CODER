package services

import (
	"fmt"
	"strings"

	entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	submissionEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
)

func makeFunctionCall(functionName string, inputs []entities.IOVariable) (string, error) {
	for _, input := range inputs {
		// For string inputs, we need to add quotes around them
		if input.Type == "string" {
			input.Value = "\"" + input.Value + "\""
		}

		// Replace input.Name in functionName with input.Value
		functionName = replaceVariable(functionName, input.Name, input.Value)
	}

	return functionName, nil
}

func isIdentifierChar(b byte) bool {
	return (b >= 'a' && b <= 'z') ||
		(b >= 'A' && b <= 'Z') ||
		(b >= '0' && b <= '9') ||
		b == '_'
}

func replaceVariable(functionCall string, variableName string, variableValue string) string {
	if variableName == "" {
		return functionCall
	}

	var out strings.Builder
	start := 0
	searchFrom := 0

	for {
		idx := strings.Index(functionCall[searchFrom:], variableName)
		if idx == -1 {
			break
		}

		idx += searchFrom
		end := idx + len(variableName)

		leftOK := idx == 0 || !isIdentifierChar(functionCall[idx-1])
		rightOK := end >= len(functionCall) || !isIdentifierChar(functionCall[end])

		if leftOK && rightOK {
			out.WriteString(functionCall[start:idx])
			out.WriteString(variableValue)
			start = end
		}

		searchFrom = end
	}

	if start == 0 {
		return functionCall
	}

	out.WriteString(functionCall[start:])
	return out.String()
}

func AppendFunctionCall(code string, function string, language submissionEntities.ProgrammingLanguage, inputs []entities.IOVariable) (string, error) {
	function_call, err := makeFunctionCall(function, inputs)
	if err != nil {
		return "", err
	}

	switch language {
	case submissionEntities.LanguagePython:
		function_call = "print(" + function_call + ")"
	case submissionEntities.LanguageJava:
		function_call = "System.out.println(" + function_call + ");"
	case submissionEntities.LanguageCPP:
		function_call = "cout << " + function_call + " << endl;"
	default:
		return "", fmt.Errorf("language not implemented")
	}

	code = code + "\n" + function_call

	return code, nil
}

func ExtractInputFromTestCase(test_case entities.TestCase) string {
	var inputs string
	for _, input := range test_case.Input {
		if inputs == "" {
			inputs = input.Value
		} else {
			inputs += "\n" + input.Value
		}
	}
	return inputs
}