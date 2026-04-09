package postarchive

import examDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"

type PathDTO struct { ID string }

func ToInput(path PathDTO) examDtos.ArchiveChallengeInput { 
	return examDtos.ArchiveChallengeInput{
		ChallengeID: path.ID,
	} 
}
