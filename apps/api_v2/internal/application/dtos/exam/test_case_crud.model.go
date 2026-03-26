package dtos

type CreateTestCaseInput struct {
	Name           string			`json:"name"`
	Input          []IOVariableDTO	`json:"input"`
	ExpectedOutput IOVariableDTO	`json:"expected_output"`
	IsSample       bool				`json:"is_sample"`
	Points         int				`json:"points"`
	ChallengeID    string			`json:"challenge_id"`
}

type UpdateTestCaseInput struct {
	ID             string			`json:"id"`
	Name           *string			`json:"name"`
	Input          *[]IOVariableDTO `json:"input"`
	ExpectedOutput *IOVariableDTO	`json:"expected_output"`
	IsSample       *bool			`json:"is_sample"`
	Points         *int				`json:"points"`
}

type GetTestCasesByChallengeInput struct {
	ChallengeID string 				`json:"challenge_id"`
	ExamID 		*string 			`json:"exam_id"`
}

type DeleteTestCaseInput struct {
	TestCaseID string 				`json:"test_case_id"`
}