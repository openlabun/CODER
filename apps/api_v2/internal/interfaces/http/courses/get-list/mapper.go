package getlist

func MapQuery(scope, studentID, teacherID string) QueryDTO {
	return QueryDTO{
		Scope: scope, 
		StudentID: studentID, 
		TeacherID: teacherID,
	}
}
