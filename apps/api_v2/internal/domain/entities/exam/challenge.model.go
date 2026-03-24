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
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Tags        []string          `json:"tags"`

	// State and Access Control
	Status       ChallengeStatus `json:"status"` // state machine: draft -> published -> archived
	Difficulty   ChallengeDifficulty `json:"difficulty"`

	// Worker Constraints
	WorkerTimeLimit    int `json:"timeLimit"` // in ms
	WorkerMemoryLimit  int `json:"memoryLimit"` // in MB
	
	// I/O Considerations
	InputVariables  []IOVariable `json:"inputVariables"`
	OutputVariable  IOVariable   `json:"outputVariable"`
	Constraints     string       `json:"constraints"`

	// Metadata
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	ExamID		    string    `json:"examId"`
	CourseID	    string    `json:"courseId"`
}
