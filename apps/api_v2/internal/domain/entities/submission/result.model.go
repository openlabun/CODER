package Submission_entities

import (
	submission_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/submission"
	ChallengeEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
)

type SubmissionResult struct {
	ID 		   string			`json:"id"`
	SubmissionID string			`json:"submission_id"`

	// Inputs for test evaluated
	TestCaseID  string			`json:"test_case_id"`

	// Result details
	Status       submission_constants.SubmissionStatus  				`json:"status"` // e.g., "accepted", "wrong_answer", "error"
	ActualOutput *ChallengeEntities.IOVariable	`json:"actual_output"` // Populated if Status is accepted/wrong_answer
	ErrorMessage *string						`json:"error_message"` // Populated if Status is error
}