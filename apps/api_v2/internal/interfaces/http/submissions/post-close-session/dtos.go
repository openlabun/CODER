package postclosesession

import submissionDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"

type PathDTO struct{ ID string }

func ToInput(path PathDTO) submissionDtos.CloseSessionInput {
	return submissionDtos.CloseSessionInput{SessionID: path.ID}
}
