package challenge_entities

import "time"

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
	ExpectedOutput *IOVariable // Optional, only populated for accepted/wrong_answer
	ActualOutput   *IOVariable // Optional, only populated for accepted/wrong_answer
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
