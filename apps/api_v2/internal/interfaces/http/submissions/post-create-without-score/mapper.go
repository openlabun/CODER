package postcreatewithoutscore

import submissionDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"

func MapRequestToInput(req RequestDTO) submissionDtos.CreateExecutionInput { return req }
