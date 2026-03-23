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
	ID 		   string
	SubmissionID string

	// Inputs for test evaluated
	TestCaseID  string

	// Result details
	Status       SubmissionStatus // e.g., "accepted", "wrong_answer", "error"
	ActualOutput *ChallengeEntities.IOVariable // Populated if Status is accepted/wrong_answer
	ErrorMessage *string // Populated if Status is error
}