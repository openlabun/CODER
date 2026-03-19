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
	Status      SubmissionStatus

	// Results
	Score       int
	TimeMsTotal int

	// Metadata
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ChallengeID string
	UserID      string
}
