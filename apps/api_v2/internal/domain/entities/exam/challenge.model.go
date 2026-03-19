package challenge_entities

import "time"

type ChallengeStatus string

const (
	ChallengeStatusDraft     ChallengeStatus = "draft"
	ChallengeStatusPublished ChallengeStatus = "published"
	ChallengeStatusArchived  ChallengeStatus = "archived"
)

type ChallengeDifficulty string

const (
	ChallengeDifficultyEasy   ChallengeDifficulty = "easy"
	ChallengeDifficultyMedium ChallengeDifficulty = "medium"
	ChallengeDifficultyHard   ChallengeDifficulty = "hard"
)


type Challenge struct {
	ID          string
	Title       string
	Description string
	Tags        []string

	// State and Access Control
	Status       ChallengeStatus // state machine: draft -> published -> archived
	Difficulty   ChallengeDifficulty

	// Worker Constraints
	WorkerTimeLimit    int // in ms
	WorkerMemoryLimit  int // in MB
	
	// I/O Considerations
	InputVariables  []IOVariable
	OutputVariable  IOVariable
	Constraints     string

	// Metadata
	CreatedAt       time.Time
	UpdatedAt       time.Time
	ExamID		    string
}
