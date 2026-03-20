package session_states

import (
	"fmt"
	"time"

	ExamEntities "../../entities/exam"
	SessionEntities "../../entities/submission"
)

// State machine for user session lifecycle
// States: active|frozen -> completed|expired|blocked
//   - active: initial state when a session starts, user can submit solutions
//   - frozen: temporary state where user cannot submit but session is not completed (e.g., due to inactivity)
//   - completed: when the user finishes the exam successfully
//   - expired: when the session time limit is exceeded without completion
//   - blocked: when the user is blocked from further attempts (e.g., due to cheating detection)

var sessionAllowedTransitions = map[SessionEntities.SessionStatus]map[SessionEntities.SessionStatus]struct{}{
	SessionEntities.SessionStatusActive: {
		SessionEntities.SessionStatusFrozen:    {},
		SessionEntities.SessionStatusCompleted: {},
		SessionEntities.SessionStatusExpired:   {},
		SessionEntities.SessionStatusBlocked:   {},
	},
	SessionEntities.SessionStatusFrozen: {
		SessionEntities.SessionStatusActive:    {},
		SessionEntities.SessionStatusCompleted: {},
		SessionEntities.SessionStatusExpired:   {},
		SessionEntities.SessionStatusBlocked:   {},
	},
}

func IsValidState(state SessionEntities.SessionStatus) bool {
	switch state {
	case SessionEntities.SessionStatusActive:
		return true
	case SessionEntities.SessionStatusFrozen:
		return true
	case SessionEntities.SessionStatusCompleted:
		return true
	case SessionEntities.SessionStatusExpired:
		return true
	case SessionEntities.SessionStatusBlocked:
		return true
	default:
		return false
	}
}

func canTransitionState(from SessionEntities.SessionStatus, to SessionEntities.SessionStatus) bool {
	if !IsValidState(from) || !IsValidState(to) {
		return false
	}

	nextStates, ok := sessionAllowedTransitions[from]
	if !ok {
		return false
	}

	_, allowed := nextStates[to]
	return allowed
}

func validateStateTransition(session SessionEntities.Session, to SessionEntities.SessionStatus) error {
	if !IsValidState(session.Status) {
		return fmt.Errorf("invalid session state: %q", session.Status)
	}

	if !IsValidState(to) {
		return fmt.Errorf("invalid target session state: %q", to)
	}

	if !canTransitionState(session.Status, to) {
		return fmt.Errorf("invalid session transition: %s -> %s", session.Status, to)
	}

	return nil
}

func ApplyTranstion(session SessionEntities.Session, to SessionEntities.SessionStatus) error {
	if err := validateStateTransition(session, to); err != nil {
		return err
	}

	session.Status = to
	return nil
}

func validateSessionExamBinding(session SessionEntities.Session, exam ExamEntities.Exam) error {
	if session.ExamID == "" {
		return fmt.Errorf("session exam id is empty")
	}

	if exam.ID == "" {
		return fmt.Errorf("exam id is empty")
	}

	if session.ExamID != exam.ID {
		return fmt.Errorf("session exam mismatch: session.ExamID=%s exam.ID=%s", session.ExamID, exam.ID)
	}

	return nil
}

func shouldExpireSession(session SessionEntities.Session, exam ExamEntities.Exam, now time.Time) bool {
	if validateStateTransition(session, SessionEntities.SessionStatusExpired) != nil {
		return true
	}

	// If exam has unlimited time, it cannot expire
	if exam.TimeLimit == -1 {
		return false
	}

	// Check if session has exceeded time limit
	if exam.TimeLimit > 0 {
		if now.After(session.StartedAt.Add(time.Duration(exam.TimeLimit) * time.Second)) {
			return true
		}
	}

	// Check if exam has ended and late submissions are not allowed
	if exam.EndTime != nil && !exam.AllowLateSubmissions {
		if now.After(*exam.EndTime) {
			return true
		}
	}

	// Check if session has no time left
	if session.TimeLeft <= 0 {
		return true
	}

	return false
}

func shouldFreezeSession(session SessionEntities.Session, now time.Time) bool {
	if validateStateTransition(session, SessionEntities.SessionStatusFrozen) != nil {
		return false
	}

	// If user has been inactive for more than 60 seconds, freeze the session
	if now.Sub(session.LastHeartbeat) > 60*time.Second { //TODO: make this configurable from env
		return true
	}

	return false
}

func shouldBlockSession(session SessionEntities.Session, exam ExamEntities.Exam) bool {
	if validateStateTransition(session, SessionEntities.SessionStatusBlocked) != nil {
		return false
	}

	// If exam has unlimited attempts, it cannot be blocked
	if exam.TryLimit == -1 {
		return false
	}

	// If user has exceeded the attempt limit, block the session
	if exam.TryLimit > 0 && session.Attempts >= exam.TryLimit {
		return true
	}

	return false
}

func UpdateSessionStatus(
	session SessionEntities.Session,
	exam ExamEntities.Exam,
	now time.Time,
	heartbeat bool,
) error {

	if err := validateSessionExamBinding(session, exam); err != nil {
		return err
	}

	// Check if session should expire before updating status
	if shouldExpireSession(session, exam, now) {
		if err := ApplyTranstion(session, SessionEntities.SessionStatusExpired); err != nil {
			return fmt.Errorf("failed to expire session: %w", err)
		}
	}

	// Check if session heartbeat is overdue (e.g., no activity for 60 seconds)
	if shouldFreezeSession(session, now) {
		if err := ApplyTranstion(session, SessionEntities.SessionStatusFrozen); err != nil {
			return fmt.Errorf("failed to freeze session due to inactivity: %w", err)
		}
	}

	// Check if session should be blocked due to attempt limit exceeded
	if shouldBlockSession(session, exam) {
		if err := ApplyTranstion(session, SessionEntities.SessionStatusBlocked); err != nil {
			return fmt.Errorf("failed to block session due to attempt limit: %w", err)
		}
	}

	if heartbeat {
		session.LastHeartbeat = now
	}

	return nil
}

func BlockSession(session *SessionEntities.Session) {
	session.Status = SessionEntities.SessionStatusBlocked
}
