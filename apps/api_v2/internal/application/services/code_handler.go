package services

import (
	submissionEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
)

func detectFunction (code string, language submissionEntities.ProgrammingLanguage, inputs []entities.IOVariable) (string, error) {
	var variable_list []string
	for _, variable := range inputs {
		variable_list = append(variable_list, variable.Name)
	}
	
	// TODO: Implement function for detecting the function name from the code.

	return "", nil
}

func makeFunctionCall (functionName string, inputs []entities.IOVariable) (string, error) {
	// TODO: Implement function for creating a function call string based on the function name and the input variables.
	
	return "", nil
}

func AppendFunctionCall (code string, language submissionEntities.ProgrammingLanguage, inputs []entities.IOVariable) (string, error) {
	function_name, err := detectFunction(code, language, inputs)
	if err != nil {
		return "", err
	}

	function_call, err := makeFunctionCall(function_name, inputs)
	if err != nil {
		return "", err
	}

	code = code + "\n" + function_call
	
	return code, nil
}
