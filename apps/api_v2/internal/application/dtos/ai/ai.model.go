package dtos

type GenerateFullChallengeInput struct {
	Topic      string `json:"topic"`
	Difficulty string `json:"difficulty"` // easy, medium, hard
}

type AIChallengeIdea struct {
	Title           string               `json:"title"`
	Description     string               `json:"description"`
	Difficulty      string               `json:"difficulty"`
	Tags            []string             `json:"tags"`
	InputFormat     string               `json:"inputFormat"`
	OutputFormat    string               `json:"outputFormat"`
	Constraints     string               `json:"constraints"`
	PublicTestCases []AITestCase `json:"publicTestCases"`
	HiddenTestCases []AITestCase `json:"hiddenTestCases"`
}

type AITestCase struct {
	Input  string `json:"input"`
	Output string `json:"output"`
	Name   string `json:"name"`
	Type   string `json:"type"` // public, hidden
}

type GenerateFullChallengeOutput struct {
	Challenge AIChallengeIdea `json:"challenge"`
}

type GenerateExamInput struct {
	Topic string `json:"topic"`
	Difficulty string `json:"difficulty"`
}

type AIExamIdea struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	TimeLimit   int    `json:"time_limit"` // in minutes
	TryLimit    int    `json:"try_limit"`
}

type GenerateExamOutput struct {
	Exam AIExamIdea `json:"exam"`
}
