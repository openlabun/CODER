package Submission_entities

import (
	ChallengeEntities "../exam"
)

type SubmissionStatus string

const (
	SubmissionStatusQueued      SubmissionStatus = "queued"
	SubmissionStatusRunning     SubmissionStatus = "running"
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