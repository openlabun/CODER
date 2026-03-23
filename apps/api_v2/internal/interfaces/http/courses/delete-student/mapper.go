package deletestudent

func MapPath(courseID, studentID string) PathDTO { 
	return PathDTO{
		CourseID: courseID, 
		StudentID: studentID,
	} 
}
