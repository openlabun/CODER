package roble_infrastructure

import (
	"context"
	"fmt"
	"strings"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
)

type SessionRepository struct {
	adapter *infrastructure.RobleDatabaseAdapter
}

func NewSessionRepository(adapter *infrastructure.RobleDatabaseAdapter) *SessionRepository {
	return &SessionRepository{adapter: adapter}
}

func (r *SessionRepository) CreateSession(ctx context.Context, session *Entities.Session) (*Entities.Session, error) {
	if session == nil {
		return nil, fmt.Errorf("session is nil")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	_, err := r.adapter.Insert(sessionTableName, []map[string]any{sessionToRecord(session)})
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (r *SessionRepository) UpdateSession(ctx context.Context, session *Entities.Session) (*Entities.Session, error) {
	if session == nil {
		return nil, fmt.Errorf("session is nil")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	sessionID := strings.TrimSpace(session.ID)
	if sessionID == "" {
		return nil, fmt.Errorf("session id is required")
	}

	_, err := r.adapter.Update(sessionTableName, "ID", sessionID, sessionToUpdates(session))
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (r *SessionRepository) DeleteSession(ctx context.Context, sessionID string) error {
	normalizedID := strings.TrimSpace(sessionID)
	if normalizedID == "" {
		return fmt.Errorf("sessionID is required")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return err
	}

	_, err := r.adapter.Delete(sessionTableName, "ID", normalizedID)
	return err
}

func (r *SessionRepository) GetSessionByID(ctx context.Context, sessionID string) (*Entities.Session, error) {
	normalizedID := strings.TrimSpace(sessionID)
	if normalizedID == "" {
		return nil, fmt.Errorf("sessionID is required")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	res, err := r.adapter.Read(sessionTableName, map[string]string{"ID": normalizedID})
	if err != nil {
		return nil, err
	}

	record, err := firstRecord(res)
	if err != nil {
		return nil, nil
	}

	return recordToSession(record)
}

func (r *SessionRepository) GetSessionsByExamID(ctx context.Context, examID string) ([]*Entities.Session, error) {
	return r.getSessionsByField(ctx, "ExamID", examID)
}

func (r *SessionRepository) GetSessionsByStudentID(ctx context.Context, studentID string) ([]*Entities.Session, error) {
	return r.getSessionsByField(ctx, "StudentID", studentID)
}

func (r *SessionRepository) getSessionsByField(ctx context.Context, field, value string) ([]*Entities.Session, error) {
	normalizedValue := strings.TrimSpace(value)
	if normalizedValue == "" {
		return nil, fmt.Errorf("%s is required", field)
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	res, err := r.adapter.Read(sessionTableName, map[string]string{field: normalizedValue})
	if err != nil {
		return nil, err
	}

	records := extractRecords(res)
	if len(records) == 0 {
		return []*Entities.Session{}, nil
	}

	sessions := make([]*Entities.Session, 0, len(records))
	for _, record := range records {
		session, mapErr := recordToSession(record)
		if mapErr != nil {
			return nil, mapErr
		}
		if session != nil {
			sessions = append(sessions, session)
		}
	}

	return sessions, nil
}
