package Submission_entities

import (
	"time"
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

type ProgrammingLanguage string

const (
	LanguagePython ProgrammingLanguage = "python"
	LanguageCPP	ProgrammingLanguage = "cpp"
	LanguageJava	ProgrammingLanguage = "java"
)

type Submission struct {
	ID          string
	Code        string
	Language    ProgrammingLanguage

	// state machine: queued -> running -> accepted|wrong_answer|error
	Status         SubmissionStatus
	ExpectedOutput *ChallengeEntities.IOVariable // Optional, only populated for accepted/wrong_answer
	ActualOutput   *ChallengeEntities.IOVariable // Optional, only populated for accepted/wrong_answer
	ErrorMessage   *string     // Optional, only populated for error state

	// Results
	Score       int
	TimeMsTotal int

	// Metadata
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ChallengeID string
	UserID      string
}
