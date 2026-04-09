package deletebyid

import examDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"

type PathDTO struct { ID string }

func ToInput(path PathDTO) examDtos.DeleteTestCaseInput { return examDtos.DeleteTestCaseInput{TestCaseID: path.ID} }
