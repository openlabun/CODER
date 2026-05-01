package getusersscoresbyexam

func MapPath(examId string) PathDTO { return PathDTO{ExamID: examId} }

func MapQuery(courseID string) QueryDTO {
	if courseID == "" {
		return QueryDTO{CourseID: nil}
	}
	return QueryDTO{CourseID: &courseID}
}
