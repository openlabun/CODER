package course_entities

import (
	"time"
	consts "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/course"
)

type Period struct {
	Year   int					`json:"year"`
	Semester consts.AcademicPeriod		`json:"semester"`
}

type Course struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"` // (Optional)

	// Access Control
	Visibility     consts.CourseVisibility `json:"visibility"`

	// Visual Identity
	VisualIdentity consts.CourseColour `json:"visual_identity"` 

	// Course Institution Data
	Code           string `json:"code"` 
	Period         *Period `json:"period"` // Optional, won't be blocked after the period ends

	// Enrollment Details
	EnrollmentCode string `json:"enrollment_code"`
	EnrollmentURL  string `json:"enrollment_url"`

	// Metadata
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	ProfessorID    string	 `json:"professor_id"`
}
