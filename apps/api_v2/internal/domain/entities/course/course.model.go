package course_entities

import "time"

type CourseColour string

const (
	CourseColourRed    CourseColour = "#ff0055"
	CourseColourGreen  CourseColour = "#00ff9d"
	CourseColourBlue   CourseColour = "#7000ff"
	CourseColourYellow CourseColour = "#ffcc00"
	CourseColourOrange  CourseColour = "#ff6b35"
	CourseColourPurple  CourseColour = "#9d00ff"
)

type CourseVisibility string

const (
	CourseVisibilityPublic  CourseVisibility = "public"
	CourseVisibilityPrivate CourseVisibility = "private"
	CourseVisibilityBlocked CourseVisibility = "blocked"
)

type AcademicPeriod string

const (
	AcademicFirstPeriod    AcademicPeriod = "01"
	AcademicIntersemestral AcademicPeriod = "02"
	AcademicSecondPeriod   AcademicPeriod = "03"
)

type Period struct {
	Year   int            `json:"year"`
	Semester AcademicPeriod `json:"semester"`
}

type Course struct {
	ID             string           `json:"id"`
	Name           string           `json:"name"`
	Description    string           `json:"description"` // (Optional)

	// Access Control
	Visibility     CourseVisibility `json:"visibility"`

	// Visual Identity
	VisualIdentity CourseColour     `json:"visual_identity"`

	// Course Institution Data
	Code           string           `json:"code"`
	Period         *Period          `json:"period"` // Optional, won't be blocked after the period ends

	// Enrollment Details
	EnrollmentCode string           `json:"enrollment_code"`
	EnrollmentURL  string           `json:"enrollment_url"`

	// Metadata
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	ProfessorID    string           `json:"professor_id"`
}
