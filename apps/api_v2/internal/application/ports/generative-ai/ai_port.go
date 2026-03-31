package generative_ai_ports

import (
	"context"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/ai"
)

type AIPort interface {
	GenerateContent(ctx context.Context, prompt string) (string, error)
	GenerateExamIdea(ctx context.Context, prompt string) (*dtos.AIExamIdea, error)
	GenerateChallengeIdea(ctx context.Context, prompt string) (*dtos.AIChallengeIdea, error)
}