package exam_factory

import (
	"strings"
	"time"
	
	"github.com/google/uuid"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	Validations "github.com/openlabun/CODER/apps/api_v2/internal/domain/validations/exam"
)


func NewChallenge(
	title, description string,
	tags []string,
	status Entities.ChallengeStatus,
	difficulty Entities.ChallengeDifficulty,
	workerTimeLimit, workerMemoryLimit int,
	inputVariables []Entities.IOVariable,
	outputVariable Entities.IOVariable,
	constraints, examID string,
) (*Entities.Challenge, error) {
	now := time.Now()
	if status == "" {
		status = Entities.ChallengeStatusDraft
	}
	if difficulty == "" {
		difficulty = Entities.ChallengeDifficultyEasy
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
		CreatedAt:         now,
		UpdatedAt:         now,
		ExamID:            strings.TrimSpace(examID),
	}

	if err := Validations.ValidateChallenge(challenge); err != nil {
		return nil, err
	}

	return challenge, nil
}

func ExistingChallenge(
	id, title, description string,
	tags []string,
	status Entities.ChallengeStatus,
	difficulty Entities.ChallengeDifficulty,
	workerTimeLimit, workerMemoryLimit int,
	inputVariables []Entities.IOVariable,
	outputVariable Entities.IOVariable,
	constraints, examID string,
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
		ExamID:            strings.TrimSpace(examID),
	}

	if err := Validations.ValidateChallenge(challenge); err != nil {
		return nil, err
	}

	return challenge, nil
}
