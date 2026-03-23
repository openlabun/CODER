package postcreate

import courseDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/course"

func MapRequestToInput(req RequestDTO) courseDtos.CreateCourseInput { 
	return req 
}
