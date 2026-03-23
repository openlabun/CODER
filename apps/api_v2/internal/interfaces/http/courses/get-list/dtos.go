package getlist

import courseDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/course"

type QueryDTO struct {
	Scope string `query:"scope"`
	StudentID string `query:"studentId"`
	TeacherID string `query:"teacherId"`
}

func ToEnrolledInput(q QueryDTO) courseDtos.GetEnrolledCoursesInput { 
	return courseDtos.GetEnrolledCoursesInput{StudentID: q.StudentID} 
}
func ToOwnedInput(q QueryDTO) courseDtos.GetOwnedCoursesInput { 
	return courseDtos.GetOwnedCoursesInput{TeacherID: q.TeacherID} 
}
