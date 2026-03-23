package mapper

import (
	factory "github.com/openlabun/CODER/apps/api_v2/internal/domain/factory/course"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"
	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/course"
)

func MapCourseStudentInputToEntity(input dtos.EnrolledInCourseInput) (*Entities.CourseStudent, error) {
	return factory.NewCourseStudent(
		input.CourseID,
		input.StudentID,
	)
}