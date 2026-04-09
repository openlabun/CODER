package postenroll

import courseDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/course"

type RequestDTO = courseDtos.EnrolledInCourseInput

func ToInput(req RequestDTO) courseDtos.EnrolledInCourseInput { 
	return req 
}
