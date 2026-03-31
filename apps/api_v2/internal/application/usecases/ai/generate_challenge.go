package ai_usecases

import (
	"context"
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/ai"
	ports "github.com/openlabun/CODER/apps/api_v2/internal/application/ports/generative-ai"
)

type GenerateFullChallengeUseCase struct {
	AI ports.AIPort
}

func NewGenerateFullChallengeUseCase(AI ports.AIPort) *GenerateFullChallengeUseCase {
	return &GenerateFullChallengeUseCase{AI: AI}
}

func (uc *GenerateFullChallengeUseCase) Execute(ctx context.Context, input dtos.GenerateFullChallengeInput) (*dtos.GenerateFullChallengeOutput, error) {
	// [STEP 1] Build prompt
	prompt := fmt.Sprintf(`
		Actúa como un diseñador experto de retos de programación competitiva (como Codeforces o LeetCode).
		Crea un reto de programación COMPLETO basado en el tema: "%s" con dificultad: "%s".

		Tu respuesta DEBE ser ÚNICAMENTE un objeto JSON válido con la siguiente estructura:
		{
		"title": "Un título creativo y descriptivo",
		"description": "Un enunciado claro y profesional del problema en español",
		"difficulty": "%s",
		"tags": ["Tag1", "Tag2"],
		"inputFormat": "Descripción detallada del formato de entrada",
		"outputFormat": "Descripción detallada del formato de salida esperado",
		"constraints": "Restricciones de tiempo y memoria (ej: N < 1000)",
		"publicTestCases": [
			{"input": "ejemplo de entrada 1", "output": "salida esperada 1", "name": "Caso Ejemplo 1", "type": "public"}
		],
		"hiddenTestCases": [
			{"input": "entrada oculta 1", "output": "salida oculta 1", "name": "Caso Evaluación 1", "type": "hidden"},
			{"input": "entrada oculta 2", "output": "salida oculta 2", "name": "Caso Evaluación 2", "type": "hidden"},
			{"input": "entrada oculta 3", "output": "salida oculta 3", "name": "Caso Evaluación 3", "type": "hidden"}
		]
		}

		REGLAS CRÍTICAS:
		1. El enunciado debe estar en ESPAÑOL.
		2. Genera EXACTAMENTE 3 casos ocultos de evaluación.
		3. El JSON debe ser perfectamente válido para ser parseado por un programa. No incluyas explicaciones fuera del JSON.
		`, input.Topic, input.Difficulty, input.Difficulty)

	// [STEP 2] Call LLM to generate challenge
	challenge, err := uc.AI.GenerateChallengeIdea(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("Error calling AI: %w", err)
	}
	if challenge == nil {
		return nil, fmt.Errorf("AI returned nil response")
	}

	return &dtos.GenerateFullChallengeOutput{Challenge: *challenge}, nil
}