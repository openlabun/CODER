package deletecourse

import courseDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/course"

type PathDTO struct {
	ID string
}

func ToInput(path PathDTO) courseDtos.DeleteCourseInput {
	return courseDtos.DeleteCourseInput{CourseID: path.ID}
}
