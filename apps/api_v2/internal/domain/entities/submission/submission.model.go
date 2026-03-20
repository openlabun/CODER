package Submission_entities

import (
	"time"
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

	// Results
	Score       int
	TimeMsTotal int

	// Metadata
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ChallengeID string
	SessionID   string
	UserID      string
}
