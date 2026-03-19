package challenge_entities

import "time"

type Visibility string

const (
	VisibilityPublic  Visibility = "public"
	VisibilityCourse  Visibility = "course"
	VisibilityTeachers Visibility = "teachers"
	VisibilityPrivate Visibility = "private"
)

type Exam struct {
	ID              string
	Title           string
	Description     string

	// Access Control
	Visibility      Visibility
	StartTime       time.Time
	EndTime         *time.Time  // Optional, null means no end time
	AllowLateSubmissions bool

	// Exam Settings
	TimeLimit int // Optional, in seconds
	TryLimit  int // Optional, null for unlimited

	// Metadata
	CreatedAt       time.Time
	UpdatedAt       time.Time
	ProfessorID     string
	CourseID        string // Optional
}
