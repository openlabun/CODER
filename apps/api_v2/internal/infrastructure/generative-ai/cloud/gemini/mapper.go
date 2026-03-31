package gemini_infrastructure

import (
	"fmt"
	"strings"
	"encoding/json"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/ai"
)

func MapResponseToChallengeIdeaDTO(geminiResponse string) (*dtos.AIChallengeIdea, error) {
	clean_json, err := cleanJSONResponse(geminiResponse)
	if err != nil {
		return nil, fmt.Errorf("Error cleaning Gemini response: %w", err)
	}
	if clean_json == "" {
		return nil, fmt.Errorf("Gemini returned empty response after cleaning. Original response: %s", geminiResponse)
	}

	var challenge dtos.AIChallengeIdea
	if err := json.Unmarshal([]byte(clean_json), &challenge); err != nil {
		return nil, fmt.Errorf("Error parsing Gemini JSON: %w. Original response: %s", err, clean_json)
	}

	return &challenge, nil
}

func MapResponseToExamIdeaDTO(geminiResponse string) (*dtos.AIExamIdea, error) {
	clean_json, err := cleanJSONResponse(geminiResponse)
	if err != nil {
		return nil, fmt.Errorf("Error cleaning Gemini response: %w", err)
	}
	if clean_json == "" {
		return nil, fmt.Errorf("Gemini returned empty response after cleaning. Original response: %s", geminiResponse)
	}

	var exam dtos.AIExamIdea
	if err := json.Unmarshal([]byte(clean_json), &exam); err != nil {
		return nil, fmt.Errorf("Error parsing Gemini JSON: %w. Original response: %s", err, clean_json)
	}

	return &exam, nil
}

func cleanJSONResponse (response string) (string, error) {
	start := strings.Index(response, "{")
	end := strings.LastIndex(response, "}")
	if start == -1 || end == -1 || end < start {
		return "", fmt.Errorf("Gemini returned invalid structure: %s\n", response)
	}

	cleanJSON := response[start : end+1]

	return cleanJSON, nil
}