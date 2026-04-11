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
	ID          string				`json:"id"`
	Code        string				`json:"code"`
	Language    ProgrammingLanguage	`json:"language"`

	// Results
	Score       int					`json:"score"`
	TimeMsTotal int					`json:"time_ms_total"`

	// Metadata
	CreatedAt   time.Time			`json:"created_at"`
	UpdatedAt   time.Time			`json:"updated_at"`
	ChallengeID string				`json:"challenge_id"`
	SessionID   string				`json:"session_id"`
	UserID      string				`json:"user_id"`
}
