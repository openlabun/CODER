package exam_usecases

import (
	"context"
	"fmt"

	sub_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/submission"
	user_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/user"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	factory "github.com/openlabun/CODER/apps/api_v2/internal/domain/factory/exam"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	mapper "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam/mapper"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
)

type GetCodeDefaultTemplates struct {
	userRepository      userRepository.UserRepository
}

func NewGetCodeDefaultTemplates(userRepository userRepository.UserRepository, examRepository examRepository.ExamRepository) *GetCodeDefaultTemplates {
	return &GetCodeDefaultTemplates{userRepository: userRepository}
}

func (uc *GetCodeDefaultTemplates) Execute(ctx context.Context, input dtos.DefaultCodeTemplatesInput) ([]Entities.CodeTemplate, error) {
	// [STEP 1] Verify user and get its role
	userEmail, err := services.UserEmailFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepository.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("user with email %q does not exist", userEmail)
	}

	if user.Role != user_constants.UserRoleProfessor {
		return nil, fmt.Errorf("user does not have permissions to access to templates")
	}

	// [STEP 2] Map Input into IOVariable entities
	inputs, output, err := mapper.MapDefaultCodeTemplatesInputToEntities(input)
	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, fmt.Errorf("output variable is required to create default template")
	}

	// [STEP 3] Create default template for each supported language
	var codeTemplates []Entities.CodeTemplate
	for _, language := range sub_constants.SupportedProgrammingLanguages {
		template, err := services.CreateTemplate(inputs, output, language)
		if err != nil {
			return nil, err
		}

		CodeTemplate, err := factory.NewCodeTemplate(string(language), template)
		if err != nil {
			return nil, err
		}
		codeTemplates = append(codeTemplates, CodeTemplate)
	}

	return codeTemplates, nil
}