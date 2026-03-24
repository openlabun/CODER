package ai_usecases

import (
	"context"
	"encoding/json"
	"fmt"
	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/ai"
	"github.com/openlabun/CODER/apps/api_v2/internal/application/services"
	"strings"
)

type GenerateExamUseCase struct {
	geminiService *services.GeminiService
}

func NewGenerateExamUseCase(geminiService *services.GeminiService) *GenerateExamUseCase {
	return &GenerateExamUseCase{geminiService: geminiService}
}

func (uc *GenerateExamUseCase) Execute(ctx context.Context, input dtos.GenerateExamInput) (*dtos.GenerateExamOutput, error) {
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

	response, err := uc.geminiService.GenerateContent(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("error llamando a Gemini: %w", err)
	}

	start := strings.Index(response, "{")
	end := strings.LastIndex(response, "}")
	if start == -1 || end == -1 || end < start {
		return nil, fmt.Errorf("la IA devolvió un formato inválido. Intenta de nuevo")
	}
	cleanJSON := response[start : end+1]

	var exam dtos.AIExamIdea
	if err := json.Unmarshal([]byte(cleanJSON), &exam); err != nil {
		return nil, fmt.Errorf("error al parsear JSON del examen")
	}

	return &dtos.GenerateExamOutput{Exam: exam}, nil
}
