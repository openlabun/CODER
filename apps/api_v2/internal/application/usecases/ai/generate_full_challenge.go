package ai_usecases

import (
	"context"
	"encoding/json"
	"fmt"
	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/ai"
	"github.com/openlabun/CODER/apps/api_v2/internal/application/services"
	"strings"
)

type GenerateFullChallengeUseCase struct {
	geminiService *services.GeminiService
}

func NewGenerateFullChallengeUseCase(geminiService *services.GeminiService) *GenerateFullChallengeUseCase {
	return &GenerateFullChallengeUseCase{geminiService: geminiService}
}

func (uc *GenerateFullChallengeUseCase) Execute(ctx context.Context, input dtos.GenerateFullChallengeInput) (*dtos.GenerateFullChallengeOutput, error) {
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

	response, err := uc.geminiService.GenerateContent(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("error llamando a Gemini: %w", err)
	}

	start := strings.Index(response, "{")
	end := strings.LastIndex(response, "}")
	if start == -1 || end == -1 || end < start {
		fmt.Printf("Gemini returned invalid structure: %s\n", response)
		return nil, fmt.Errorf("la IA devolvió un formato inválido. Intenta de nuevo")
	}
	cleanJSON := response[start : end+1]

	var challenge dtos.AIChallengeIdea
	if err := json.Unmarshal([]byte(cleanJSON), &challenge); err != nil {
		fmt.Printf("Error al parsear JSON de Gemini: %v\nRespuesta original: %s\n", err, cleanJSON)
		return nil, fmt.Errorf("la IA devolvió un formato inválido. Prueba con un tema más sencillo")
	}

	return &dtos.GenerateFullChallengeOutput{Challenge: challenge}, nil
}
