package dtos

import (
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
)

type CreateSubmissionInput struct {
	Code        string		`json:"code"`
	Language    string		`json:"language"`
	Score       int			`json:"score"`
	ChallengeID string		`json:"challenge_id"`
	SessionID   string		`json:"session_id"`
}

type UpdateResultInput struct {
	ResultID	string		`json:"result_id"`
	Status		string		`json:"status"`
	TimeExecution	int		`json:"time_execution"`
	Output		*string		`json:"output"`
	Error		*string		`json:"error"`
}

type GetSubmissionStatusInput struct {
	SubmissionID string		`json:"submission_id"`
}

type GetUserSubmissionsInput struct {
	UserID      string		`json:"user_id"`
	Status      *string		`json:"status"`			// Optional
	TestID      *string		`json:"test_id"`		// Optional
	ChallengeID *string 	`json:"challenge_id"`	// Optional
}

type GetChallengeSubmissionsInput struct {
	ChallengeID string		`json:"challenge_id"`
	Status      *string		`json:"status"`			// Optional
	TestID      *string		`json:"test_id"`		// Optional
}

type GetSessionSubmissionsInput struct {
	SessionID   string		`json:"session_id"`
	Status      *string		`json:"status"`			// Optional
	TestID      *string		`json:"test_id"`		// Optional
	ChallengeID *string		`json:"challenge_id"`	// Optional
}

type SubmissionOutputDTO struct {
	Submission Entities.Submission			`json:"submission"`
	Results    []Entities.SubmissionResult	`json:"results"`
}

type SubmissionResultPublishedDTO struct {
	SubmissionID string `json:"submission_id"`
	Code 	  	 string `json:"code"`
	ResultID	 string `json:"result_id"`
	TimeLimitMs  int    `json:"time_limit_ms"`
	MemoryLimitMb int    `json:"memory_limit_mb"`
	Status		 string `json:"status"`
	Type 		 string `json:"type"`
	Language 	 string `json:"language"`
}
