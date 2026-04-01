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
	ID           string		 `json:"id"`
	StudentID    string		 `json:"user_id"`
	ExamID       string		 `json:"exam_id"`

	// Status
	// State machine: active|frozen -> completed|expired|blocked
	Status 	 SessionStatus	 `json:"status"`
	Attempts int			 `json:"attempts"`
	TimeLeft int 			 `json:"time_left"` // in seconds, -1 for unlimited
	
	// Metadata
	StartedAt	  time.Time	 `json:"started_at"`
	FinishedAt    *time.Time `json:"finished_at"` // Optional, null means session is still active
	LastHeartbeat time.Time	 `json:"last_heartbeat"`
}