package Submission_entities

import (
	"time"

	submission_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/submission"
)

type Session struct {
	ID           string		 `json:"id"`
	StudentID    string		 `json:"user_id"`
	ExamID       string		 `json:"exam_id"`

	// Status
	// State machine: active|frozen -> completed|expired|blocked
	Status 	 submission_constants.SessionStatus	 `json:"status"`
	Attempts int			 `json:"attempts"`
	TimeLeft int 			 `json:"time_left"` // in seconds, -1 for unlimited
	
	// Metadata
	StartedAt	  time.Time	 `json:"started_at"`
	FinishedAt    *time.Time `json:"finished_at"` // Optional, null means session is still active
	LastHeartbeat time.Time	 `json:"last_heartbeat"`
}