package postfork

import examDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"

type PathDTO struct {
	ID string `params:"id"`
}

func ToInput(p PathDTO) examDtos.ForkChallengeInput {
	return examDtos.ForkChallengeInput{
		ChallengeID: p.ID,
	}
}
