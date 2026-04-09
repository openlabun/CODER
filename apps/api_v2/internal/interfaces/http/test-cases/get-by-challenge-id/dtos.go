package getbychallengeid

import examDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"

type PathDTO struct { ChallengeID string }

func ToInput(path PathDTO) examDtos.GetTestCasesByChallengeInput { return examDtos.GetTestCasesByChallengeInput{ChallengeID: path.ChallengeID} }
