package cache_publisher

import (
	"context"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	ExamRepositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	cache_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
	PublisherAdapter "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/cache-handler/publisher"
)

type ChallengeRepository struct {
	publisherAdapter *PublisherAdapter.CachePublisherAdapter
	repository       ExamRepositories.ChallengeRepository
}

func NewChallengeRepository(publisherAdapter *PublisherAdapter.CachePublisherAdapter, repository ExamRepositories.ChallengeRepository) *ChallengeRepository {
	return &ChallengeRepository{publisherAdapter: publisherAdapter, repository: repository}
}

func (r *ChallengeRepository) CreateChallenge(ctx context.Context, challenge *Entities.Challenge) (*Entities.Challenge, error) {
	message, err := PublisherAdapter.MapChallengeToSyncMessage(challenge, cache_infrastructure.InsertOperation)
	if err != nil {
		return nil, err
	}
	if err := r.publisherAdapter.Save(message); err != nil {
		return nil, err
	}
	return challenge, nil
}

func (r *ChallengeRepository) UpdateChallenge(ctx context.Context, challenge *Entities.Challenge) (*Entities.Challenge, error) {
	message, err := PublisherAdapter.MapChallengeToSyncMessage(challenge, cache_infrastructure.UpdateOperation)
	if err != nil {
		return nil, err
	}
	if err := r.publisherAdapter.Save(message); err != nil {
		return nil, err
	}
	return challenge, nil
}

func (r *ChallengeRepository) DeleteChallenge(ctx context.Context, challengeID string) error {
	return r.repository.DeleteChallenge(ctx, challengeID)
}

func (r *ChallengeRepository) GetChallenges(ctx context.Context, status, tag, difficulty *string) ([]*Entities.Challenge, error) {
	return r.repository.GetChallenges(ctx, status, tag, difficulty)
}

func (r *ChallengeRepository) GetChallengeByID(ctx context.Context, challengeID string) (*Entities.Challenge, error) {
	return r.repository.GetChallengeByID(ctx, challengeID)
}

func (r *ChallengeRepository) GetChallengesByExamID(ctx context.Context, examID string) ([]*Entities.Challenge, error) {
	return r.repository.GetChallengesByExamID(ctx, examID)
}

func (r *ChallengeRepository) GetChallengesByUserID(ctx context.Context, userID string, examID *string) ([]*Entities.Challenge, error) {
	return r.repository.GetChallengesByUserID(ctx, userID, examID)
}

func (r *ChallengeRepository) GetChallengesByTag(ctx context.Context, tag string) ([]*Entities.Challenge, error) {
	return r.repository.GetChallengesByTag(ctx, tag)
}

func (r *ChallengeRepository) GetInputVariablesByChallengeID(ctx context.Context, challengeID string) ([]*Entities.IOVariable, error) {
	return r.repository.GetInputVariablesByChallengeID(ctx, challengeID)
}

func (r *ChallengeRepository) GetOutputVariablesByChallengeID(ctx context.Context, challengeID string) ([]*Entities.IOVariable, error) {
	return r.repository.GetOutputVariablesByChallengeID(ctx, challengeID)
}
