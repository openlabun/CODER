package challenge_entities

import (
	"time"

	consts "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/exam"
)

type Exam struct {
	ID              string					`json:"id"`
	Title           string					`json:"title"`
	Description     string					`json:"description"`

	// Access Control
	Visibility      consts.Visibility		`json:"visibility"` // public, course, teachers, private
	StartTime       time.Time				`json:"start_time"`
	EndTime         *time.Time				`json:"end_time"`  // Optional, null means no end time
	AllowLateSubmissions bool				`json:"allow_late_submissions"`

	// Exam Settings
	TimeLimit int							`json:"time_limit"` // Optional, in seconds, -1 for unlimited
	TryLimit  int							`json:"try_limit"`  // Optional, -1 for unlimited

	// Metadata
	CreatedAt       time.Time				`json:"created_at"`
	UpdatedAt       time.Time				`json:"updated_at"`
	ProfessorID     string					`json:"professor_id"`
	CourseID        *string					`json:"course_id"` // Optional
}

func (e *Exam) IsOpen() bool {
	now := time.Now()
	return (e.StartTime.Before(now) && (e.EndTime == nil || e.EndTime.After(now))) || e.AllowLateSubmissions
}