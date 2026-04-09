package getlist

func MapQuery(examID string) QueryDTO {
	var ptr *string
	if examID != "" {
		ptr = &examID
	}
	return QueryDTO{
		ExamID: ptr,
	}
}
