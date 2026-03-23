package postpublish

import examDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"

type PathDTO struct { ID string }

func ToInput(path PathDTO) examDtos.PublishChallengeInput { 
	return examDtos.PublishChallengeInput{
		ChallengeID: path.ID,
	} 
}
