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
	ID              string      `json:"id"`
	Title           string      `json:"title"`
	Description     string      `json:"description"`

	// Access Control
	Visibility      Visibility  `json:"visibility"`
	StartTime       time.Time   `json:"startTime"`
	EndTime         *time.Time  `json:"endTime"`  // Optional, null means no end time
	AllowLateSubmissions bool   `json:"allowLateSubmissions"`

	// Exam Settings
	TimeLimit int `json:"timeLimit"` // Optional, in seconds, -1 for unlimited
	TryLimit  int `json:"tryLimit"`  // Optional, -1 for unlimited

	// Metadata
	CreatedAt       time.Time   `json:"createdAt"`
	UpdatedAt       time.Time   `json:"updatedAt"`
	ProfessorID     string      `json:"professorId"`
	CourseID        string      `json:"courseId"` // Optional
}
