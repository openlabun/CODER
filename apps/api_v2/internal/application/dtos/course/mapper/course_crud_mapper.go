package mapper

import (
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"
	factory "github.com/openlabun/CODER/apps/api_v2/internal/domain/factory/course"
	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/course"
	consts "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/course"
)

func MapCreateCourseInputToCourseEntity(professorID string, input dtos.CreateCourseInput) (*Entities.Course, error) {
	// Map fields from CreateCourseInput to Course entity
	period := Entities.Period{
		Year: input.Year,
		Semester: consts.AcademicPeriod(input.Semester),
	}

	visibility := consts.CourseVisibility(input.Visibility)
	visual_identity := consts.CourseColour(input.VisualIdentity)

	return factory.NewCourse(
		input.Name,
		input.Description,
		visibility,
		visual_identity,
		input.Code,
		&period,
		input.EnrollmentCode,
		input.EnrollmentURL,
		professorID,
	)
}

func MapUpdateCourseInputToCourseEntity(originalCourse *Entities.Course, input dtos.UpdateCourseInput) (*Entities.Course, error) {
	// Map fields from UpdateCourseInput to Course entity
	if input.Name != nil {
		originalCourse.Name = *input.Name
	}
	if input.Description != nil {
		originalCourse.Description = *input.Description
	}
	if input.Visibility != nil {
		originalCourse.Visibility = consts.CourseVisibility(*input.Visibility)
	}
	if input.VisualIdentity != nil {
		originalCourse.VisualIdentity = consts.CourseColour(*input.VisualIdentity)
	}
	if input.Code != nil {
		originalCourse.Code = *input.Code
	}
	if input.Year != nil && input.Semester != nil {
		originalCourse.Period = &Entities.Period{
			Year: *input.Year,
			Semester: consts.AcademicPeriod(*input.Semester),
		}
	}
	if input.EnrollmentCode != nil {
		originalCourse.EnrollmentCode = *input.EnrollmentCode
	}
	if input.EnrollmentURL != nil {
		originalCourse.EnrollmentURL = *input.EnrollmentURL
	}

	return originalCourse, nil
}