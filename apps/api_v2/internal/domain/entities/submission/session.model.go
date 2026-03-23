package Submission_entities

import (
	"time"
)

type SessionStatus string

const (
	SessionStatusActive   SessionStatus = "active"
	SessionStatusCompleted SessionStatus = "completed"
	SessionStatusExpired   SessionStatus = "expired"
	SessionStatusBlocked   SessionStatus = "blocked"
	SessionStatusFrozen	SessionStatus = "frozen"
)

type Session struct {
	ID           string
	StudentID    string
	ExamID       string

	// Status
	// State machine: active|frozen -> completed|expired|blocked
	Status 	 SessionStatus
	Attempts int
	TimeLeft int // in seconds, -1 for unlimited
	
	// Metadata
	StartedAt	time.Time
	FinishedAt  *time.Time // Optional, null means session is still active
	LastHeartbeat time.Time
}