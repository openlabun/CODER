package challenge_entities

import (
	"time"

	exam_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/exam"
)

type ChallengeStatus = exam_constants.ChallengeStatus

const (
	ChallengeStatusDraft     ChallengeStatus = exam_constants.ChallengeStatusDraft
	ChallengeStatusPublished ChallengeStatus = exam_constants.ChallengeStatusPublished
	ChallengeStatusPrivate   ChallengeStatus = exam_constants.ChallengeStatusPrivate
	ChallengeStatusArchived  ChallengeStatus = exam_constants.ChallengeStatusArchived
)

type ChallengeDifficulty = exam_constants.ChallengeDifficulty

const (
	ChallengeDifficultyEasy   ChallengeDifficulty = exam_constants.ChallengeDifficultyEasy
	ChallengeDifficultyMedium ChallengeDifficulty = exam_constants.ChallengeDifficultyMedium
	ChallengeDifficultyHard   ChallengeDifficulty = exam_constants.ChallengeDifficultyHard
)

type Challenge struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`

	// State and Access Control
	Status     ChallengeStatus     `json:"status"` // state machine: draft -> published|private -> archived
	Difficulty ChallengeDifficulty `json:"difficulty"`

	// Worker Constraints
	WorkerTimeLimit   int `json:"worker_time_limit"`   // in ms
	WorkerMemoryLimit int `json:"worker_memory_limit"` // in MB

	// I/O Considerations
	InputVariables []IOVariable `json:"input_variables"`
	OutputVariable IOVariable   `json:"output_variable"`
	Constraints    string       `json:"constraints"`

	// Metadata
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    string    `json:"user_id"`
}
