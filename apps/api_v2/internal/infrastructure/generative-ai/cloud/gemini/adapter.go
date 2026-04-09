package gemini_infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/ai"
)

type GeminiAdapter struct {
	url    string
	apiKey string
}

func NewGeminiAdapter() *GeminiAdapter {
	return &GeminiAdapter{
		url:   "https://generativelanguage.googleapis.com/v1beta/models/gemini-flash-latest:generateContent",
		apiKey: os.Getenv("GEMINI_API_KEY"),
	}
}

func (s *GeminiAdapter) GenerateContent(ctx context.Context, prompt string) (string, error) {
	if s.apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY is not set")
	}

	body := map[string]any{
		"contents": []map[string]any{
			{
				"parts": []map[string]any{
					{"text": prompt},
				},
			},
		},
		"generationConfig": map[string]any{
			"response_mime_type": "application/json",
		},
	}

	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, "POST", s.url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", s.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		fmt.Printf("Gemini API Error Body: %+v\n", errResp)
		return "", fmt.Errorf("gemini api returned status %d", resp.StatusCode)
	}

	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from gemini")
	}

	return result.Candidates[0].Content.Parts[0].Text, nil
}

func (s *GeminiAdapter) GenerateExamIdea(ctx context.Context, prompt string) (*dtos.AIExamIdea, error) {
	response, err := s.GenerateContent(ctx, prompt)
	if err != nil {
		return nil, err
	}
	return MapResponseToExamIdeaDTO(response)
}

func (s *GeminiAdapter) GenerateChallengeIdea(ctx context.Context, prompt string) (*dtos.AIChallengeIdea, error) {
	response, err := s.GenerateContent(ctx, prompt)
	if err != nil {
		return nil, err
	}
	return MapResponseToChallengeIdeaDTO(response)
}

