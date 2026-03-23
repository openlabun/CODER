package deletebyid

import challengeDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"

type PathDTO struct{ ID string }

func ToInput(path PathDTO) challengeDtos.DeleteChallengeInput {
	return challengeDtos.DeleteChallengeInput{ChallengeID: path.ID}
}
