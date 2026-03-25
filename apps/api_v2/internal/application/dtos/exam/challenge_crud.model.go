package dtos

type IOVariableDTO struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type CreateChallengeInput struct {
	Title             string            `json:"title"`
	Description       string            `json:"description"`
	Tags              []string          `json:"tags"`
	Status            string            `json:"status"`
	Difficulty        string            `json:"difficulty"`
	WorkerTimeLimit   int               `json:"workerTimeLimit"`
	WorkerMemoryLimit int               `json:"workerMemoryLimit"`
	InputVariables    []IOVariableDTO  `json:"inputVariables"`
	OutputVariable    IOVariableDTO    `json:"outputVariable"`
	Constraints       string            `json:"constraints"`
	CreatedAt         string            `json:"createdAt"`
	UpdatedAt         string            `json:"updatedAt"`
	ExamID            string            `json:"examId"`
}

type UpdateChallengeInput struct {
	ChallengeID       string             `json:"challenge_id"`
	Title             *string           `json:"title"`
	Description       *string           `json:"description"`
	Tags              *[]string          `json:"tags"`
	Status            *string           `json:"status"`
	Difficulty        *string           `json:"difficulty"`
	WorkerTimeLimit   *int              `json:"workerTimeLimit"`
	WorkerMemoryLimit *int              `json:"workerMemoryLimit"`
	InputVariables    *[]IOVariableDTO   `json:"inputVariables"`
	OutputVariable    *IOVariableDTO     `json:"outputVariable"`
	Constraints       *string           `json:"constraints"`
	CreatedAt         *string           `json:"createdAt"`
	UpdatedAt         *string           `json:"updatedAt"`
	ExamID            *string           `json:"examId"`
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