package roble_infrastructure

import (
	"context"
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

func (r *SubmissionRepository) CreateSubmission(ctx context.Context, submission *Entities.Submission) (*Entities.Submission, error) {
	if submission == nil {
		return nil, fmt.Errorf("submission is nil")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	_, err := r.adapter.Insert(submissionTableName, []map[string]any{submissionToRecord(submission)})
	if err != nil {
		return nil, err
	}

	return submission, nil
}

func (r *SubmissionRepository) UpdateSubmission(ctx context.Context, submission *Entities.Submission) (*Entities.Submission, error) {
	if submission == nil {
		return nil, fmt.Errorf("submission is nil")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
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

func (r *SubmissionRepository) DeleteSubmission(ctx context.Context, submissionID string) error {
	normalizedID := strings.TrimSpace(submissionID)
	if normalizedID == "" {
		return fmt.Errorf("submissionID is required")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return err
	}

	_, err := r.adapter.Delete(submissionTableName, "ID", normalizedID)
	return err
}

func (r *SubmissionRepository) GetSubmissionByID(ctx context.Context, submissionID string) (*Entities.Submission, error) {
	normalizedID := strings.TrimSpace(submissionID)
	if normalizedID == "" {
		return nil, fmt.Errorf("submissionID is required")
	}

	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
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

func (r *SubmissionRepository) GetSubmissionsByExamItemScoreID(ctx context.Context, examItemScoreID string) ([]*Entities.Submission, error) {
	normalizedExamItemScoreID := strings.TrimSpace(examItemScoreID)
	if normalizedExamItemScoreID == "" {
		return nil, fmt.Errorf("examItemScoreID is required")
	}

	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	res, err := r.adapter.Read(submissionTableName, map[string]string{"ExamItemScoreID": normalizedExamItemScoreID})
	if err != nil {
		return nil, err
	}

	records := extractRecords(res)
	if len(records) == 0 {
		return nil, nil
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


func (r *SubmissionRepository) GetLastSubmissionByExamItemScoreID(ctx context.Context, examItemScoreID string) (*Entities.Submission, error) {
	submission, err := r.GetSubmissionsByExamItemScoreID(ctx, examItemScoreID)
	if err != nil {
		return nil, err
	}

	var lastestSubmission *Entities.Submission
	for _, submission := range submission {
		if submission == nil {
			return nil, nil
		}

		if lastestSubmission == nil {
			lastestSubmission = submission
			continue
		}
		
		if submission.UpdatedAt.After(lastestSubmission.UpdatedAt) {
			lastestSubmission = submission
		}
	}

	return lastestSubmission, nil
}

func (r *SubmissionRepository) GetBestSubmissionByExamItemScoreID(ctx context.Context, examItemScoreID string) (*Entities.Submission, error) {
	submission, err := r.GetSubmissionsByExamItemScoreID(ctx, examItemScoreID)
	if err != nil {
		return nil, err
	}

	var bestSubmission *Entities.Submission
	for _, submission := range submission {
		if submission == nil {
			return nil, nil
		}

		if bestSubmission == nil {
			bestSubmission = submission
			continue
		}

		if submission.Score > bestSubmission.Score {
			bestSubmission = submission
		}
	}

	return bestSubmission, nil
}

func (r *SubmissionRepository) GetSubmissionsBySessionID(ctx context.Context, sessionID string, status *string, testID *string, challengeID *string) ([]*Entities.Submission, error) {
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	conditions := r.buildConditions(nil, status, testID, challengeID, &sessionID)
	return r.getSubmissionsByFields(conditions)
}

func (r *SubmissionRepository) GetSubmissionsByUserID(ctx context.Context, userID string, status *string, testID *string, challengeID *string) ([]*Entities.Submission, error) {
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	conditions := r.buildConditions(&userID, status, testID, challengeID, nil)
	return r.getSubmissionsByFields(conditions)
}

func (r *SubmissionRepository) GetSubmissionsByChallengeID(ctx context.Context, challengeID string, status *string, testID *string) ([]*Entities.Submission, error) {
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	conditions := r.buildConditions(nil, status, testID, &challengeID, nil)
	return r.getSubmissionsByFields(conditions)
}

func (r *SubmissionRepository) buildConditions(userID *string, status *string, testID *string, challengeID *string, sessionID *string) map[string]string {
	conditions := map[string]string{}
	if status != nil {
		conditions["Status"] = *status
	}
	if testID != nil {
		conditions["TestID"] = *testID
	}
	if challengeID != nil {
		conditions["ChallengeID"] = *challengeID
	}
	if userID != nil {
		conditions["UserID"] = *userID
	}

	if sessionID != nil {
		conditions["SessionID"] = *sessionID
	}

	return conditions
}

func (r *SubmissionRepository) getSubmissionsByFields(fields map[string]string) ([]*Entities.Submission, error) {

	res, err := r.adapter.Read(submissionTableName, fields)
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
