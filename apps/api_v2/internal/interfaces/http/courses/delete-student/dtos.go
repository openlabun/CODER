package deletestudent

import courseDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/course"

type PathDTO struct { CourseID string; StudentID string }

func ToInput(path PathDTO) courseDtos.RemoveStudentFromCourseInput {
	return courseDtos.RemoveStudentFromCourseInput{
		CourseID: path.CourseID, 
		StudentID: path.StudentID,
	}
}
