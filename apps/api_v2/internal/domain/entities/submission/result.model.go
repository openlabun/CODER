package Submission_entities

import (
	ChallengeEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
)

type SubmissionStatus string

const (
	SubmissionStatusQueued      SubmissionStatus = "queued"
	SubmissionStatusRunning     SubmissionStatus = "running"
	SubmissionStatusTimeout     SubmissionStatus = "timeout"
	SubmissionStatusExecuted	SubmissionStatus = "executed"
	SubmissionStatusAccepted    SubmissionStatus = "accepted"
	SubmissionStatusWrongAnswer SubmissionStatus = "wrong_answer"
	SubmissionStatusError       SubmissionStatus = "error"
)

type SubmissionResult struct {
	ID 		   string			`json:"id"`
	SubmissionID string			`json:"submission_id"`

	// Inputs for test evaluated
	TestCaseID  string			`json:"test_case_id"`

	// Result details
	Status       SubmissionStatus  				`json:"status"` // e.g., "accepted", "wrong_answer", "error"
	ActualOutput *ChallengeEntities.IOVariable	`json:"actual_output"` // Populated if Status is accepted/wrong_answer
	ErrorMessage *string						`json:"error_message"` // Populated if Status is error
}