package getstudents

import courseDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/course"

type PathDTO struct { CourseID string }

func ToInput(path PathDTO) courseDtos.GetCourseStudentsInput { 
	return courseDtos.GetCourseStudentsInput{CourseID: path.CourseID} 
}
