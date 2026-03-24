package challenge_entities

import "time"

type TestCase struct {
	ID             string `json:"id"`
	Name           string

	// I/O Configuration
	Input          []IOVariable
	ExpectedOutput IOVariable

	// Scoring
	IsSample       bool
	Points         int

	// Metadata
	CreatedAt      time.Time
	ChallengeID    string
}
