package postblock

import submissionDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"

type PathDTO struct{ ID string }

func ToInput(path PathDTO) submissionDtos.BlockSessionInput {
	return submissionDtos.BlockSessionInput{SessionID: path.ID}
}
