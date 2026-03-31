package postgeneratefullchallenge

import aiDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/ai"

func MapRequestToInput(r RequestDTO) aiDtos.GenerateFullChallengeInput {
	return aiDtos.GenerateFullChallengeInput{
		Topic:      r.Topic,
		Difficulty: r.Difficulty,
	}
}
