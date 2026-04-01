package postgenerateexam

import aiDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/ai"

func MapRequestToInput(r RequestDTO) aiDtos.GenerateExamInput {
	return aiDtos.GenerateExamInput{
		Topic:      r.Topic,
		Difficulty: r.Difficulty,
	}
}
