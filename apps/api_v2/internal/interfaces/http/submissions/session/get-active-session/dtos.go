package getbysessionid

import submissionDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"

type PathDTO struct {
	UserID *string
}

func ToInput(path PathDTO) submissionDtos.GetActiveSessionInput {
	return submissionDtos.GetActiveSessionInput{UserID: path.UserID}
}
