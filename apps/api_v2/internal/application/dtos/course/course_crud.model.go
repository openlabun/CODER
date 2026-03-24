package dtos

type CreateCourseInput struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	Visibility     string `json:"visibility"`
	VisualIdentity string `json:"visual_identity"`
	Code           string `json:"code"`
	Year           int    `json:"year"`
	Semester       string `json:"semester"`
	EnrollmentCode string `json:"enrollment_code"`
	EnrollmentURL  string `json:"enrollment_url"`
	TeacherID      string `json:"teacher_id"`
}

type DeleteCourseInput struct {
	CourseID string `json:"course_id"`
}

type UpdateCourseInput struct {
	ID 		   	   string `json:"id"`
	Name           *string `json:"name"`
	Description    *string `json:"description"`
	Visibility     *string `json:"visibility"`
	VisualIdentity *string `json:"visual_identity"`
	Code           *string `json:"code"`
	Year           *int    `json:"year"`
	Semester       *string `json:"semester"`
	EnrollmentCode *string `json:"enrollment_code"`
	EnrollmentURL  *string `json:"enrollment_url"`
	TeacherID      *string `json:"teacher_id"`
}

type GetCourseDetailsInput struct {
	CourseID string `json:"course_id"`
}

type GetEnrolledCoursesInput struct {
	StudentID string `json:"student_id"`
}

type GetOwnedCoursesInput struct {
	TeacherID string `query:"teacherId"`
}

type CourseBrowseItem struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Code           string `json:"code"`
	Period         string `json:"period"`
	ProfessorID    string `json:"professor_id"`
	ProfessorName  string `json:"professor_name"`
	CreatedAt      string `json:"created_at"`
}