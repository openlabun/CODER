package getexamitems

func MapQuery(examID string) QueryDTO {
	return QueryDTO{
		ExamID: examID,
	}
}
