package postaddstudent

import courseDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/course"

type PathDTO struct { 
	CourseID string 
}

type RequestDTO struct { 
	StudentID    *string `json:"studentID"` 
	StudentEmail *string `json:"studentEmail"`
}

func ToInput(path PathDTO, body RequestDTO) courseDtos.EnrolledInCourseInput {
	return courseDtos.EnrolledInCourseInput{
		CourseID: path.CourseID, 
		StudentID: body.StudentID,
		StudentEmail: body.StudentEmail,
	}
}
