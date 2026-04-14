package Submission_entities

import (
	"time"

	submission_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/submission"
)

type Submission struct {
	ID          	 string				`json:"id"`
	Code        	 string				`json:"code"`
	Language    	 submission_constants.ProgrammingLanguage	`json:"language"`

	// Results
	Score       	 int				`json:"score"`
	TimeMsTotal 	 int				`json:"time_ms_total"`
	Scorable		 bool 				`json:"scorable"`

	// Metadata
	CreatedAt   	 time.Time			`json:"created_at"`
	UpdatedAt   	 time.Time			`json:"updated_at"`
	ChallengeID 	 string				`json:"challenge_id"`
	SessionID   	 string				`json:"session_id"`
	UserID     		 string				`json:"user_id"`
	ExamItemScoreID  string				`json:"exam_item_score_id"`
}
