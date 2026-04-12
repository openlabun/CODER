package dtos

type IOVariableDTO struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type CodeTemplateDTO struct {
	Language string `json:"language"`
	Template string `json:"template"`
}

type CreateChallengeInput struct {
	Title             string            `json:"title"`
	Description       string            `json:"description"`
	Tags              []string          `json:"tags"`
	Status            string            `json:"status"`
	Difficulty        string            `json:"difficulty"`
	WorkerTimeLimit   int               `json:"worker_time_limit"`
	WorkerMemoryLimit int               `json:"worker_memory_limit"`
	CodeTemplates     []CodeTemplateDTO `json:"code_templates"`
	InputVariables    []IOVariableDTO   `json:"input_variables"`
	OutputVariable    IOVariableDTO     `json:"output_variable"`
	Constraints       string            `json:"constraints"`
	UserID            string            `json:"user_id"`
}

type UpdateChallengeInput struct {
	ChallengeID       string             `json:"challenge_id"`
	Title             *string            `json:"title"`
	Description       *string            `json:"description"`
	Tags              *[]string          `json:"tags"`
	Status            *string            `json:"status"`
	Difficulty        *string            `json:"difficulty"`
	WorkerTimeLimit   *int               `json:"worker_time_limit"`
	WorkerMemoryLimit *int               `json:"worker_memory_limit"`
	CodeTemplates     *[]CodeTemplateDTO `json:"code_templates"`
	InputVariables    *[]IOVariableDTO   `json:"input_variables"`
	OutputVariable    *IOVariableDTO     `json:"output_variable"`
	Constraints       *string            `json:"constraints"`
	UserID            *string            `json:"user_id"`
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

type GetChallengesByUserInput struct {
	ExamID *string `json:"exam_id"`
}

type GetPublicChallengesInput struct {
	Tag        *string `json:"tag"`
	Difficulty *string `json:"difficulty"`
}

type ForkChallengeInput struct {
	ChallengeID string `json:"challenge_id"`
}
