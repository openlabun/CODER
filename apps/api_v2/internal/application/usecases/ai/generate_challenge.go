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

		Tu respuesta DEBE ser ÚNICAMENTE un objeto JSON válido con la siguiente estructura y formato estricto:
		{
		"title": "...",
		"description": "...",
		"difficulty": "%s",
		"tags": ["..."],
		"worker_time_limit": 1500,
		"worker_memory_limit": 256,
		"input_variables": [
			{"name": "nombre_variable", "type": "tipo", "value": ""}
		],
		"output_variable": {"name": "nombre_resultado", "type": "tipo", "value": ""},
		"constraints": "...",
		"public_test_cases": [
			{
				"name": "Ejemplo 1",
				"type": "public",
				"input": [
					{"name": "nombre_variable", "type": "tipo", "value": "valor_ejemplo"}
				],
				"output": {"name": "nombre_resultado", "type": "tipo", "value": "resultado_ejemplo"}
			}
		],
		"hidden_test_cases": [
			{
				"name": "Oculto 1",
				"type": "hidden",
				"input": [
					{"name": "nombre_variable", "type": "tipo", "value": "valor_oculto"}
				],
				"output": {"name": "nombre_resultado", "type": "tipo", "value": "resultado_oculto"}
			}
		]
		}

		REGLAS CRÍTICAS DE DISEÑO:
		1. FIDELIDAD AL TEMA: Debes seguir estrictamente el "Tema del Reto" proporcionado. Si el tema es simple (ej: "Suma dos números"), el reto debe ser directo y fiel a esa lógica básica. NO intentes transformar un tema simple en un problema complejo de algoritmos famosos (como Two Sum) usando historias espaciales o fantasiosas innecesarias. Sé literal y profesional.
		2. DIFICULTAD REALISTA: Respeta el nivel "%s". 
		   - 'easy': El problema debe ser directo, ideal para principiantes, sin algoritmos complejos.
		   - 'medium': Requiere lógica moderada, uso de arreglos, estructuras de datos básicas o algoritmos simples.
		   - 'hard': Requiere algoritmos avanzados, optimización de tiempo/memoria, estructuras de datos complejas o múltiples pasos lógicos.
		3. LENGUAJE: Todo el enunciado y los nombres de los casos de prueba deben estar en ESPAÑOL.
		4. VARIABLES DINÁMICAS: Crea variables con nombres que tengan sentido para el problema. El campo 'type' puede ser 'string', 'int', 'float', 'boolean', 'array'. 
		5. ESTRUCTURA DE TEST CASES: El campo 'input' DEBE coincidir con las variables definidas. El campo 'value' DEBE ser SIEMPRE un string.
		6. CALIDAD: Genera exactamente 3 hidden_test_cases y 2 public_test_cases con valores de entrada variados y lógicos.
		7. JSON ESTRICTO: Tu respuesta debe ser un JSON puro, sin explicaciones ni código markdown.
		`, input.Topic, input.Difficulty, input.Difficulty, input.Difficulty)

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