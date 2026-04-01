package getlist

type QueryDTO struct {
	Scope string `query:"scope"`
	StudentID string `query:"studentId"`
	TeacherID string `query:"teacherId"`
}

