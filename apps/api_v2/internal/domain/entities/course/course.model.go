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
	AcademicFirstPeriod   AcademicPeriod = "01"
	AcademicIntersemestral AcademicPeriod = "02"
	AcademicSecondPeriod AcademicPeriod = "03"
)

type Period struct {
	Year   int
	Semester AcademicPeriod
}

type Course struct {
	ID             string
	Name           string
	Description    string // (Optional)

	// Access Control
	Visibility     CourseVisibility

	// Visual Identity
	VisualIdentity CourseColour

	// Course Institution Data
	Code           string
	Period         Period

	// Enrollment Details
	EnrollmentCode string
	EnrollmentURL  string

	// Metadata
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ProfessorID    string
}
