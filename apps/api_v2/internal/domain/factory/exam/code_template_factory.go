package exam_factory

import (
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/submission"
)

func NewCodeTemplate (language string, template string) (Entities.CodeTemplate, error) {
	return Entities.CodeTemplate{
		Language: constants.ProgrammingLanguage(language),
		Template: template,
	}, nil
}

func ExistingCodeTemplate (language string, template string) (Entities.CodeTemplate, error) {
	return Entities.CodeTemplate{
		Language: constants.ProgrammingLanguage(language),
		Template: template,
	}, nil
}
