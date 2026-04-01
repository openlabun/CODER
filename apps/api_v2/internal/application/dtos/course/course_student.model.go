package dtos

type EnrolledInCourseInput struct {
	CourseID  	 string  `json:"course_id"`
	StudentID 	 *string `json:"student_id"`
	StudentEmail *string `json:"student_email"`
}

type GetCourseStudentsInput struct {
	CourseID string `json:"course_id"`
}

type RemoveStudentFromCourseInput struct {
	CourseID  string `json:"course_id"`
	StudentID string `json:"student_id"`
}