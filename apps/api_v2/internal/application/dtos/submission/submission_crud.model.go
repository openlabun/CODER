package dtos

import (
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
)

type CreateSubmissionInput struct {
	Code        string
	Function    string
	Language    string
	Score       int
	ChallengeID string
	SessionID   string
}

type UpdateResultInput struct {
	ResultID	string
	Status		string
	TimeExecution	int
	Output		*string
	Error		*string
}

type GetSubmissionStatusInput struct {
	SubmissionID string
}

type GetUserSubmissionsInput struct {
	UserID      string
	Status      *string
	TestID      *string
	ChallengeID *string
}

type GetChallengeSubmissionsInput struct {
	ChallengeID string
	Status      *string
	TestID      *string
}

type SubmissionOutputDTO struct {
	Submission Entities.Submission
	Results    []Entities.SubmissionResult
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
