package dtos

type CreateTestCaseInput struct {
	Name           string
	Input          []IOVariableDTO
	ExpectedOutput IOVariableDTO
	IsSample       bool
	Points         int
	ChallengeID    string
}

type UpdateTestCaseInput struct {
	ID             string
	Name           *string
	Input          *[]IOVariableDTO
	ExpectedOutput *IOVariableDTO
	IsSample       *bool
	Points         *int
}

type GetTestCasesByChallengeInput struct {
	ChallengeID string `json:"challenge_id"`
}

type DeleteTestCaseInput struct {
	TestCaseID string `json:"test_case_id"`
}