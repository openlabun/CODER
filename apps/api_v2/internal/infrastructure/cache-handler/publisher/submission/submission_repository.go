package cache_publisher

import (
	"context"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	SubmissionRepositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
	cache_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
	PublisherAdapter "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/cache-handler/publisher"
)

type SubmissionRepository struct {
	publisherAdapter *PublisherAdapter.CachePublisherAdapter
	repository       SubmissionRepositories.SubmissionRepository
}

func NewSubmissionRepository(publisherAdapter *PublisherAdapter.CachePublisherAdapter, repository SubmissionRepositories.SubmissionRepository) *SubmissionRepository {
	return &SubmissionRepository{publisherAdapter: publisherAdapter, repository: repository}
}

func (r *SubmissionRepository) CreateSubmission(ctx context.Context, submission *Entities.Submission) (*Entities.Submission, error) {
	message, err := PublisherAdapter.MapSubmissionToSyncMessage(submission, cache_infrastructure.InsertOperation)
	if err != nil {
		return nil, err
	}
	if err := r.publisherAdapter.Save(message); err != nil {
		return nil, err
	}
	return submission, nil
}

func (r *SubmissionRepository) UpdateSubmission(ctx context.Context, submission *Entities.Submission) (*Entities.Submission, error) {
	message, err := PublisherAdapter.MapSubmissionToSyncMessage(submission, cache_infrastructure.UpdateOperation)
	if err != nil {
		return nil, err
	}
	if err := r.publisherAdapter.Save(message); err != nil {
		return nil, err
	}
	return submission, nil
}

func (r *SubmissionRepository) DeleteSubmission(ctx context.Context, submissionID string) error {
	return r.repository.DeleteSubmission(ctx, submissionID)
}

func (r *SubmissionRepository) GetSubmissionByID(ctx context.Context, submissionID string) (*Entities.Submission, error) {
	return r.repository.GetSubmissionByID(ctx, submissionID)
}

func (r *SubmissionRepository) GetSubmissionsByExamItemScoreID(ctx context.Context, examItemScoreID string) ([]*Entities.Submission, error) {
	return r.repository.GetSubmissionsByExamItemScoreID(ctx, examItemScoreID)
}

func (r *SubmissionRepository) GetBestSubmissionByExamItemScoreID(ctx context.Context, examItemScoreID string) (*Entities.Submission, error) {
	return r.repository.GetBestSubmissionByExamItemScoreID(ctx, examItemScoreID)
}

func (r *SubmissionRepository) GetLastSubmissionByExamItemScoreID(ctx context.Context, examItemScoreID string) (*Entities.Submission, error) {
	return r.repository.GetLastSubmissionByExamItemScoreID(ctx, examItemScoreID)
}

func (r *SubmissionRepository) GetSubmissionsBySessionID(ctx context.Context, sessionID string, status *string, testID *string, challengeID *string) ([]*Entities.Submission, error) {
	return r.repository.GetSubmissionsBySessionID(ctx, sessionID, status, testID, challengeID)
}

func (r *SubmissionRepository) GetSubmissionsByUserID(ctx context.Context, userID string, status *string, testID *string, challengeID *string) ([]*Entities.Submission, error) {
	return r.repository.GetSubmissionsByUserID(ctx, userID, status, testID, challengeID)
}

func (r *SubmissionRepository) GetSubmissionsByChallengeID(ctx context.Context, challengeID string, status *string, testID *string) ([]*Entities.Submission, error) {
	return r.repository.GetSubmissionsByChallengeID(ctx, challengeID, status, testID)
}
