package postsession

import submissionDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"

type RequestDTO = submissionDtos.CreateSessionInput

func ToInput(req RequestDTO) submissionDtos.CreateSessionInput { return req }
