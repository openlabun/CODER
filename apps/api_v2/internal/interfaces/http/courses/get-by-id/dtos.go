package getbyid

import courseDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/course"

type PathDTO struct { ID string }

func ToInput(path PathDTO) courseDtos.GetCourseDetailsInput { 
	return courseDtos.GetCourseDetailsInput{CourseID: path.ID} 
}
