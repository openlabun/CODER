package challenge_entities

import "time"

type TestCase struct {
	ID             string		`json:"id"`
	Name           string		`json:"name"`

	// I/O Configuration
	Input          []IOVariable	`json:"input"`
	ExpectedOutput IOVariable	`json:"expected_output"`

	// Scoring
	IsSample       bool			`json:"is_sample"`
	Points         int			`json:"points"`

	// Metadata
	CreatedAt      time.Time	`json:"created_at"`
	ChallengeID    string		`json:"challenge_id"`
}
