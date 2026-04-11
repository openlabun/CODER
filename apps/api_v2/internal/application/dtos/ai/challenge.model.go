package dtos

import examDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"

type GenerateFullChallengeInput struct {
	Topic      string `json:"topic"`
	Difficulty string `json:"difficulty"` // easy, medium, hard
}

type AIChallengeIdea struct {
	Title             string                   `json:"title"`
	Description       string                   `json:"description"`
	Difficulty        string                   `json:"difficulty"`
	Tags              []string                 `json:"tags"`
	InputVariables    []examDtos.IOVariableDTO `json:"input_variables"`
	OutputVariable    examDtos.IOVariableDTO   `json:"output_variable"`
	WorkerTimeLimit   int                      `json:"worker_time_limit"`
	WorkerMemoryLimit int                      `json:"worker_memory_limit"`
	Constraints       string                   `json:"constraints"`
	PublicTestCases   []AITestCase             `json:"public_test_cases"`
	HiddenTestCases   []AITestCase             `json:"hidden_test_cases"`
}

type AITestCase struct {
	Input   []examDtos.IOVariableDTO `json:"input"`
	Output  examDtos.IOVariableDTO   `json:"output"`
	Name    string                   `json:"name"`
	Visible bool                     `json:"visible"` // If visible is true, it's a public testcase
}

type GenerateFullChallengeOutput struct {
	Challenge AIChallengeIdea `json:"challenge"`
}