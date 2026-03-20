package submission_factory

import (
	"strings"
	"time"

	Entities "../../entities/submission"
	Validations "../../validations/submission"
)

func NewSubmission(id, code string, language Entities.ProgrammingLanguage, challengeID, sessionID, userID string) (*Entities.Submission, error) {
	now := time.Now()
	submission := &Entities.Submission{
		ID:          strings.TrimSpace(id),
		Code:        code,
		Language:    language,
		Score:       0,
		TimeMsTotal: 0,
		CreatedAt:   now,
		UpdatedAt:   now,
		ChallengeID: strings.TrimSpace(challengeID),
		SessionID:   strings.TrimSpace(sessionID),
		UserID:      strings.TrimSpace(userID),
	}

	if err := Validations.ValidateSubmission(submission); err != nil {
		return nil, err
	}

	return submission, nil
}

func ExistingSubmission(
	id, code string,
	language Entities.ProgrammingLanguage,
	score, timeMsTotal int,
	createdAt, updatedAt time.Time,
	challengeID, sessionID, userID string,
) (*Entities.Submission, error) {
	submission := &Entities.Submission{
		ID:          strings.TrimSpace(id),
		Code:        code,
		Language:    language,
		Score:       score,
		TimeMsTotal: timeMsTotal,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		ChallengeID: strings.TrimSpace(challengeID),
		SessionID:   strings.TrimSpace(sessionID),
		UserID:      strings.TrimSpace(userID),
	}

	if err := Validations.ValidateSubmission(submission); err != nil {
		return nil, err
	}

	return submission, nil
}

