package ai_usecases

import (
	"context"
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/ai"
	ports "github.com/openlabun/CODER/apps/api_v2/internal/application/ports/generative-ai"
)

type GenerateExamUseCase struct {
	AI ports.AIPort
}

func NewGenerateExamUseCase(AI ports.AIPort) *GenerateExamUseCase {
	return &GenerateExamUseCase{AI: AI}
}

func (uc *GenerateExamUseCase) Execute(ctx context.Context, input dtos.GenerateExamInput) (*dtos.GenerateExamOutput, error) {
	// [STEP 1] Build prompt for LLM
	prompt := fmt.Sprintf(`
		Actúa como un profesor universitario experto diseñando exámenes de programación.
		Crea un examen completo basado en el tema: "%s" con dificultad general: "%s".

		Tu respuesta DEBE ser ÚNICAMENTE un objeto JSON válido con la siguiente estructura:
		{
		"title": "Un título profesional para el examen",
		"description": "Una descripción clara que incluya instrucciones y reglas generales en español",
		"time_limit": 90,
		"try_limit": 2
		}

		REGLAS CRÍTICAS:
		1. El contenido debe estar en ESPAÑOL.
		2. time_limit debe ser un entero representando minutos (entre 30 y 180).
		3. try_limit debe ser un entero pequeño (entre 1 y 5).
		4. El JSON debe ser perfectamente válido. No incluyas explicaciones fuera del JSON.
	`, input.Topic, input.Difficulty)

	// [STEP 2] Call LLM to generate exam
	exam, err := uc.AI.GenerateExamIdea(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("Error calling AI: %w", err)
	}
	if exam == nil {
		return nil, fmt.Errorf("AI returned nil response")
	}

	return &dtos.GenerateExamOutput{Exam: *exam}, nil
}