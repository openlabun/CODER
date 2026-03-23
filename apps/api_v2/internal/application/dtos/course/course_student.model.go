package dtos

type EnrolledInCourseInput struct {
	CourseID  string `json:"course_id"`
	StudentID string `json:"student_id"`
}

type GetCourseStudentsInput struct {
	CourseID string `json:"course_id"`
}

type RemoveStudentFromCourseInput struct {
	CourseID  string `json:"course_id"`
	StudentID string `json:"student_id"`
}