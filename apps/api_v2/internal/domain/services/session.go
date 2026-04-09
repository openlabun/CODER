package services

import (
	"fmt"
	"context"
	"time"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	ExamEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	submissionRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
	session_states "github.com/openlabun/CODER/apps/api_v2/internal/domain/states/session"
)

func GetUpdatedSession (ctx context.Context, sessionID string, exam *ExamEntities.Exam, now time.Time, sessionRepository submissionRepository.SessionRepository) (*Entities.Session, error) {
	session, err := sessionRepository.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if session == nil {
		return nil, fmt.Errorf("session not found")
	}

	// Check if session should be frozen or expired based on last heartbeat
	if err := session_states.UpdateSessionStatus(session, exam, now, false); err != nil {
		return nil, err
	}

	return session, nil
}