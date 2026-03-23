package getbycourseid

import examDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"

type RequestDTO struct { CourseID string }

func (r RequestDTO) ToInput() examDtos.GetExamsByCourseInput { 
	return examDtos.GetExamsByCourseInput{CourseID: r.CourseID} 
}
