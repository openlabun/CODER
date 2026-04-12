package exam_factory

import (
	"strings"
	"time"

	"github.com/google/uuid"

	exam_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/exam"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	Validations "github.com/openlabun/CODER/apps/api_v2/internal/domain/validations/exam"
)


func NewChallenge(
	title, description string,
	tags []string,
	status exam_constants.ChallengeStatus,
	difficulty exam_constants.ChallengeDifficulty,
	workerTimeLimit, workerMemoryLimit int,
	inputVariables []Entities.IOVariable,
	outputVariable Entities.IOVariable,
	constraints, UserID string,
) (*Entities.Challenge, error) {
	now := time.Now()
	if status == "" {
		status = exam_constants.ChallengeStatusDraft
	}
	if difficulty == "" {
		difficulty = exam_constants.ChallengeDifficultyEasy
	}

	challenge := &Entities.Challenge{
		ID:                uuid.New().String(),
		Title:             strings.TrimSpace(title),
		Description:       strings.TrimSpace(description),
		Tags:              tags,
		Status:            status,
		Difficulty:        difficulty,
		WorkerTimeLimit:   workerTimeLimit,
		WorkerMemoryLimit: workerMemoryLimit,
		InputVariables:    inputVariables,
		OutputVariable:    outputVariable,
		Constraints:       strings.TrimSpace(constraints),
		UserID:            strings.TrimSpace(UserID),
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err := Validations.ValidateChallenge(challenge); err != nil {
		return nil, err
	}

	return challenge, nil
}

func ExistingChallenge(
	id, title, description string,
	tags []string,
	status exam_constants.ChallengeStatus,
	difficulty exam_constants.ChallengeDifficulty,
	workerTimeLimit, workerMemoryLimit int,
	inputVariables []Entities.IOVariable,
	outputVariable Entities.IOVariable,
	constraints, UserID string,
	createdAt, updatedAt time.Time,
) (*Entities.Challenge, error) {
	challenge := &Entities.Challenge{
		ID:                strings.TrimSpace(id),
		Title:             strings.TrimSpace(title),
		Description:       strings.TrimSpace(description),
		Tags:              tags,
		Status:            status,
		Difficulty:        difficulty,
		WorkerTimeLimit:   workerTimeLimit,
		WorkerMemoryLimit: workerMemoryLimit,
		InputVariables:    inputVariables,
		OutputVariable:    outputVariable,
		Constraints:       strings.TrimSpace(constraints),
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
		UserID:            strings.TrimSpace(UserID),
	}

	if err := Validations.ValidateChallenge(challenge); err != nil {
		return nil, err
	}

	return challenge, nil
}
