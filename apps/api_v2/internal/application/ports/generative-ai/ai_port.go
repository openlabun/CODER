package generative_ai_ports

import (
	"context"
)

type AIPort interface {
	GenerateContent(ctx context.Context, prompt string) (string, error)
}