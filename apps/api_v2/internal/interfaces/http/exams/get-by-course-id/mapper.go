package getbycourseid

func MapPathToRequest(courseID string) RequestDTO { 
	return RequestDTO{
		CourseID: courseID,
	} 
}
