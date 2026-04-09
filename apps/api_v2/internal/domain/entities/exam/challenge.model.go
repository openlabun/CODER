package challenge_entities

import "time"

type ChallengeStatus string

const (
	ChallengeStatusDraft     ChallengeStatus = "draft"
	ChallengeStatusPublished ChallengeStatus = "published"
	ChallengeStatusPrivate  ChallengeStatus = "private"
	ChallengeStatusArchived  ChallengeStatus = "archived"
)

type ChallengeDifficulty string

const (
	ChallengeDifficultyEasy   ChallengeDifficulty = "easy"
	ChallengeDifficultyMedium ChallengeDifficulty = "medium"
	ChallengeDifficultyHard   ChallengeDifficulty = "hard"
)


type Challenge struct {
	ID          string	 				`json:"id"`
	Title       string					`json:"title"`
	Description string					`json:"description"`
	Tags        []string				`json:"tags"`

	// State and Access Control
	Status       ChallengeStatus        `json:"status"` // state machine: draft -> published|private -> archived
	Difficulty   ChallengeDifficulty    `json:"difficulty"`

	// Worker Constraints
	WorkerTimeLimit    int  			`json:"worker_time_limit"` // in ms
	WorkerMemoryLimit  int 				`json:"worker_memory_limit"` // in MB
	
	// I/O Considerations
	InputVariables  []IOVariable		`json:"input_variables"`
	OutputVariable  IOVariable			`json:"output_variable"`
	Constraints     string				`json:"constraints"`

	// Metadata
	CreatedAt       time.Time			`json:"created_at"`
	UpdatedAt       time.Time			`json:"updated_at"`
	UserID		  	string				`json:"user_id"`
}
