package roble_infrastructure

import (
	"fmt"
	"strings"
	"time"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
)

type SubmissionRepository struct {
	adapter *infrastructure.RobleDatabaseAdapter
}

func NewSubmissionRepository(adapter *infrastructure.RobleDatabaseAdapter) *SubmissionRepository {
	return &SubmissionRepository{adapter: adapter}
}

func (r *SubmissionRepository) CreateSubmission(submission *Entities.Submission) (*Entities.Submission, error) {
	if submission == nil {
		return nil, fmt.Errorf("submission is nil")
	}

	_, err := r.adapter.Insert(submissionTableName, []map[string]any{submissionToRecord(submission)})
	if err != nil {
		return nil, err
	}

	return submission, nil
}

func (r *SubmissionRepository) UpdateSubmission(submission *Entities.Submission) (*Entities.Submission, error) {
	if submission == nil {
		return nil, fmt.Errorf("submission is nil")
	}

	submissionID := strings.TrimSpace(submission.ID)
	if submissionID == "" {
		return nil, fmt.Errorf("submission id is required")
	}

	updates := submissionToUpdates(submission)
	updates["UpdatedAt"] = time.Now().UTC().Format(time.RFC3339)

	_, err := r.adapter.Update(submissionTableName, "ID", submissionID, updates)
	if err != nil {
		return nil, err
	}

	submission.ID = submissionID
	if ts, ok := updates["UpdatedAt"].(string); ok {
		if parsed, parseErr := time.Parse(time.RFC3339, ts); parseErr == nil {
			submission.UpdatedAt = parsed
		}
	}

	return submission, nil
}

func (r *SubmissionRepository) DeleteSubmission(submissionID string) error {
	normalizedID := strings.TrimSpace(submissionID)
	if normalizedID == "" {
		return fmt.Errorf("submissionID is required")
	}

	_, err := r.adapter.Delete(submissionTableName, "ID", normalizedID)
	return err
}

func (r *SubmissionRepository) GetSubmissionByID(submissionID string) (*Entities.Submission, error) {
	normalizedID := strings.TrimSpace(submissionID)
	if normalizedID == "" {
		return nil, fmt.Errorf("submissionID is required")
	}

	res, err := r.adapter.Read(submissionTableName, map[string]string{"ID": normalizedID})
	if err != nil {
		return nil, err
	}

	record, err := firstRecord(res)
	if err != nil {
		return nil, nil
	}

	return recordToSubmission(record)
}

func (r *SubmissionRepository) GetSubmissionsBySessionID(sessionID string) ([]*Entities.Submission, error) {
	return r.getSubmissionsByField("SessionID", sessionID)
}

func (r *SubmissionRepository) GetSubmissionsByUserID(userID string) ([]*Entities.Submission, error) {
	return r.getSubmissionsByField("UserID", userID)
}

func (r *SubmissionRepository) GetSubmissionsByChallengeID(challengeID string) ([]*Entities.Submission, error) {
	return r.getSubmissionsByField("ChallengeID", challengeID)
}

func (r *SubmissionRepository) getSubmissionsByField(field, value string) ([]*Entities.Submission, error) {
	normalizedValue := strings.TrimSpace(value)
	if normalizedValue == "" {
		return nil, fmt.Errorf("%s is required", field)
	}

	res, err := r.adapter.Read(submissionTableName, map[string]string{field: normalizedValue})
	if err != nil {
		return nil, err
	}

	records := extractRecords(res)
	if len(records) == 0 {
		return []*Entities.Submission{}, nil
	}

	submissions := make([]*Entities.Submission, 0, len(records))
	for _, record := range records {
		submission, mapErr := recordToSubmission(record)
		if mapErr != nil {
			return nil, mapErr
		}
		if submission != nil {
			submissions = append(submissions, submission)
		}
	}

	return submissions, nil
}
