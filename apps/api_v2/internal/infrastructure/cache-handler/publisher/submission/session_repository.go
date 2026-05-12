package cache_publisher

import (
	"context"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	SubmissionRepositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
	cache_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
	PublisherAdapter "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/cache-handler/publisher"
)

type SessionRepository struct {
	publisherAdapter *PublisherAdapter.CachePublisherAdapter
	repository       SubmissionRepositories.SessionRepository
}

func NewSessionRepository(publisherAdapter *PublisherAdapter.CachePublisherAdapter, repository SubmissionRepositories.SessionRepository) *SessionRepository {
	return &SessionRepository{publisherAdapter: publisherAdapter, repository: repository}
}

func (r *SessionRepository) CreateSession(ctx context.Context, session *Entities.Session) (*Entities.Session, error) {
	message, err := PublisherAdapter.MapSessionToSyncMessage(session, cache_infrastructure.InsertOperation)
	if err != nil {
		return nil, err
	}
	if err := r.publisherAdapter.Save(message); err != nil {
		return nil, err
	}
	return session, nil
}

func (r *SessionRepository) UpdateSession(ctx context.Context, session *Entities.Session) (*Entities.Session, error) {
	message, err := PublisherAdapter.MapSessionToSyncMessage(session, cache_infrastructure.UpdateOperation)
	if err != nil {
		return nil, err
	}
	if err := r.publisherAdapter.Save(message); err != nil {
		return nil, err
	}
	return session, nil
}

func (r *SessionRepository) DeleteSession(ctx context.Context, sessionID string) error {
	return r.repository.DeleteSession(ctx, sessionID)
}

func (r *SessionRepository) GetSessionByID(ctx context.Context, sessionID string) (*Entities.Session, error) {
	return r.repository.GetSessionByID(ctx, sessionID)
}

func (r *SessionRepository) GetSessionsByExamID(ctx context.Context, examID string) ([]*Entities.Session, error) {
	return r.repository.GetSessionsByExamID(ctx, examID)
}

func (r *SessionRepository) GetSessionsByStudentID(ctx context.Context, studentID string) ([]*Entities.Session, error) {
	return r.repository.GetSessionsByStudentID(ctx, studentID)
}
