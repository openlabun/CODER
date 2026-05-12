package cache_roble_infrastructure

import (
	"context"
	"strings"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	cache "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
	roble_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
	repository "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/submission"
)

type SessionRepository struct {
	robleRepository *repository.SessionRepository
	cacheAdapter    *cache.CacheAdapter
}

func NewSessionRepository(robleAdapter *roble_infrastructure.RobleDatabaseAdapter, cacheAdapter *cache.CacheAdapter) *SessionRepository {
	return &SessionRepository{
		robleRepository: repository.NewSessionRepository(robleAdapter),
		cacheAdapter:    cacheAdapter,
	}
}

func (r *SessionRepository) CreateSession(ctx context.Context, session *Entities.Session) (*Entities.Session, error) {
	result, err := r.robleRepository.CreateSession(ctx, session)
	if err != nil {
		return nil, err
	}
	if rec, e := cache.EntityToMap(result); e == nil {
		r.cacheAdapter.Save(cache.CacheSessionTable,  cache.InsertOperation, rec)
	}
	return result, nil
}

func (r *SessionRepository) UpdateSession(ctx context.Context, session *Entities.Session) (*Entities.Session, error) {
	result, err := r.robleRepository.UpdateSession(ctx, session)
	if err != nil {
		return nil, err
	}
	if rec, e := cache.EntityToMap(result); e == nil {
		r.cacheAdapter.Save(cache.CacheSessionTable, cache.UpdateOperation, rec)
	}
	return result, nil
}

func (r *SessionRepository) DeleteSession(ctx context.Context, sessionID string) error {
	if err := r.robleRepository.DeleteSession(ctx, sessionID); err != nil {
		return err
	}
	r.cacheAdapter.DeleteByID(cache.CacheSessionTable, strings.TrimSpace(sessionID))
	return nil
}

func (r *SessionRepository) GetSessionByID(ctx context.Context, sessionID string) (*Entities.Session, error) {
	id := strings.TrimSpace(sessionID)
	if rec, err := r.cacheAdapter.FindByID(cache.CacheSessionTable, id); err == nil && rec != nil {
		if s, e := cache.MapToEntity[Entities.Session](rec); e == nil {
			return s, nil
		}
	}
	s, err := r.robleRepository.GetSessionByID(ctx, sessionID)
	if err != nil || s == nil {
		return s, err
	}
	if rec, e := cache.EntityToMap(s); e == nil {
		r.cacheAdapter.Save(cache.CacheSessionTable,  cache.InsertOperation, rec)
	}
	return s, nil
}

func (r *SessionRepository) GetSessionsByExamID(ctx context.Context, examID string) ([]*Entities.Session, error) {
	id := strings.TrimSpace(examID)
	if recs, err := r.cacheAdapter.FindBy(cache.CacheSessionTable, map[string]string{"exam_id": id}); err == nil && recs != nil {
		return sessionsFromRecords(recs), nil
	}
	sessions, err := r.robleRepository.GetSessionsByExamID(ctx, examID)
	if err != nil {
		return nil, err
	}
	for _, s := range sessions {
		if rec, e := cache.EntityToMap(s); e == nil {
			r.cacheAdapter.Save(cache.CacheSessionTable,  cache.InsertOperation, rec)
		}
	}
	return sessions, nil
}

func (r *SessionRepository) GetSessionsByStudentID(ctx context.Context, studentID string) ([]*Entities.Session, error) {
	id := strings.TrimSpace(studentID)
	if recs, err := r.cacheAdapter.FindBy(cache.CacheSessionTable, map[string]string{"user_id": id}); err == nil && recs != nil {
		return sessionsFromRecords(recs), nil
	}
	sessions, err := r.robleRepository.GetSessionsByStudentID(ctx, studentID)
	if err != nil {
		return nil, err
	}
	for _, s := range sessions {
		if rec, e := cache.EntityToMap(s); e == nil {
			r.cacheAdapter.Save(cache.CacheSessionTable,  cache.InsertOperation, rec)
		}
	}
	return sessions, nil
}