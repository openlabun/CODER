package postcreate

import examDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"

func MapRequestToInput(req RequestDTO) examDtos.CreateExamInput { return req }
