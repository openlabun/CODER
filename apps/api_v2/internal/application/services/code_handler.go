package services

import (
	"fmt"

	exam_consts "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/exam"
	constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/submission"
	entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
)

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

func CreateTemplate (inputs []entities.IOVariable, output *entities.IOVariable, language constants.ProgrammingLanguage) (string, error) {
	var template string
	// Get if language is now available
	if !languageIsSupported(language) {
		return "", fmt.Errorf("language %q is not supported for default template creation", language)
	}

	switch language {
		case constants.LanguagePython:
			template = buildTemplatePython(inputs, output)
		case constants.LanguageJava:
			template = "//NOT IMPLEMENTED"
		case constants.LanguageCPP:
			template = "//NOT IMPLEMENTED"
		default:
			return "", fmt.Errorf("Language not supported for default template building")
	}

	return template, nil
}

func languageIsSupported(language constants.ProgrammingLanguage) bool {
	for _, supportedLanguage := range constants.SupportedProgrammingLanguages {
		if language == supportedLanguage {
			return true
		}
	}
	return false
}

func buildTemplatePython (inputs []entities.IOVariable, output *entities.IOVariable) string {
	template := "import sys\n\n"
	if len(inputs) > 0 {
		template += "Variables = sys.stdin.read().split()\n\n"
	}
	template += createOutputDeclarationPython(output) + "\n"
	template += createInputsCallPython(inputs)
	template += "\n# Write your code here\n\n"
	template += createOutputPrintPython(output)
	return template
}

func createInputsCallPython (inputs []entities.IOVariable) string {
	inputsCall := ""
	for i, input := range inputs {
		// Append a line of (var_name) = VarType(Variables[i]) for each input
		switch input.Type {
			case exam_consts.VariableFormatInt:
				inputsCall += fmt.Sprintf("%s = int(Variables[%d])\n", input.Name, i)
			case exam_consts.VariableFormatFloat:
				inputsCall += fmt.Sprintf("%s = float(Variables[%d])\n", input.Name, i)
			case exam_consts.VariableFormatString:
				inputsCall += fmt.Sprintf("%s = str(Variables[%d])\n", input.Name, i)
		}
	}

	return inputsCall
}  

func createOutputPrintPython (output *entities.IOVariable) string {
	outputPrint := ""
	if output == nil {
		return outputPrint
	}

	switch output.Type {
		case exam_consts.VariableFormatInt, exam_consts.VariableFormatFloat, exam_consts.VariableFormatString:
			outputPrint = fmt.Sprintf("print(%s)\n", output.Name)
	}
	return outputPrint
}

func createOutputDeclarationPython (output *entities.IOVariable) string {
	outputDeclaration := ""
	if output == nil {
		return outputDeclaration
	}

	switch output.Type {
		case exam_consts.VariableFormatInt:
			outputDeclaration = fmt.Sprintf("%s = 0\n", output.Name)
		case exam_consts.VariableFormatFloat:
			outputDeclaration = fmt.Sprintf("%s = 0.0\n", output.Name)
		case exam_consts.VariableFormatString:
			outputDeclaration = fmt.Sprintf("%s = \"\"\n", output.Name)
	}
	return outputDeclaration
}