package getpublic

import examDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"

type QueryDTO struct {
	Tag        *string `query:"tag"`
	Difficulty *string `query:"difficulty"`
}

func ToInput(q QueryDTO) examDtos.GetPublicChallengesInput {
	return examDtos.GetPublicChallengesInput{
		Tag:        q.Tag,
		Difficulty: q.Difficulty,
	}
}
