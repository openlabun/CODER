package submission_validations

import (
	"fmt"
	"strings"

	StateMachine "github.com/openlabun/CODER/apps/api_v2/internal/domain/states/session"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
)

func ValidateSession(session *Entities.Session) error {
	if session == nil {
		return fmt.Errorf("session is nil")
	}

	if strings.TrimSpace(session.ID) == "" {
		return fmt.Errorf("session id is required")
	}

	if strings.TrimSpace(session.StudentID) == "" {
		return fmt.Errorf("session student id is required")
	}

	if strings.TrimSpace(session.ExamID) == "" {
		return fmt.Errorf("session exam id is required")
	}

	if !StateMachine.IsValidState(session.Status) {
		return fmt.Errorf("invalid session status: %q", session.Status)
	}

	if session.Attempts < 0 {
		return fmt.Errorf("session attempts cannot be negative")
	}

	if session.TimeLeft < 0 {
		return fmt.Errorf("session time left cannot be negative")
	}

	if session.StartedAt.IsZero() {
		return fmt.Errorf("session startedAt is required")
	}

	if session.LastHeartbeat.IsZero() {
		return fmt.Errorf("session lastHeartbeat is required")
	}

	if session.FinishedAt != nil && session.FinishedAt.Before(session.StartedAt) {
		return fmt.Errorf("session finishedAt cannot be before startedAt")
	}

	return nil
}
