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
	InputFormat     string               `json:"input_format"`
	OutputFormat    string               `json:"output_format"`
	Constraints     string               `json:"constraints"`
	PublicTestCases []AITestCase 		 `json:"public_test_cases"`
	HiddenTestCases []AITestCase 		 `json:"hidden_test_cases"`
}

type AITestCase struct {
	Input  string `json:"input"`
	Output string `json:"output"`
	Name   string `json:"name"`
	Visible   bool `json:"visible"` // If visible is true, it's a public testcase
}

type GenerateFullChallengeOutput struct {
	Challenge AIChallengeIdea `json:"challenge"`
}