package dtos

type IOVariableDTO struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type CreateChallengeInput struct {
	Title       string
	Description string
	Tags        []string
	Status       string
	Difficulty   string
	WorkerTimeLimit    int
	WorkerMemoryLimit  int
	InputVariables  []IOVariableDTO
	OutputVariable  IOVariableDTO
	Constraints     string
	CreatedAt       string
	UpdatedAt       string
	ExamID		    string
}

type UpdateChallengeInput struct {
	ChallengeID string `json:"challenge_id"`
	Title       *string
	Description *string
	Tags        *[]string
	Status       *string
	Difficulty   *string
	WorkerTimeLimit    *int
	WorkerMemoryLimit  *int
	InputVariables  *[]IOVariableDTO
	OutputVariable  *IOVariableDTO
	Constraints     *string
	CreatedAt       *string
	UpdatedAt       *string
	ExamID		    *string
}

type DeleteChallengeInput struct {
	ChallengeID string `json:"challenge_id"`
}

type PublishChallengeInput struct {
	ChallengeID string `json:"challenge_id"`
}

type ArchiveChallengeInput struct {
	ChallengeID string `json:"challenge_id"`
}

type GetChallengeDetailsInput struct {
	ChallengeID string `json:"challenge_id"`
}

type GetChallengesByExamInput struct {
	ExamID string `json:"exam_id"`
}