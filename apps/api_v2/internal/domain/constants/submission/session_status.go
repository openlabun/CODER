package constants

type SessionStatus string

const (
	SessionStatusActive    SessionStatus = "active"
	SessionStatusCompleted SessionStatus = "completed"
	SessionStatusExpired   SessionStatus = "expired"
	SessionStatusBlocked   SessionStatus = "blocked"
	SessionStatusFrozen    SessionStatus = "frozen"
)
