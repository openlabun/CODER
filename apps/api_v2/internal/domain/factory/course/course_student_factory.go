package course_factory

import (
	"strings"

	Entities "../../entities/course"
	Validations "../../validations/course"
)

func NewCourseStudent(courseID, studentID string) (*Entities.CourseStudent, error) {
	courseStudent := &Entities.CourseStudent{
		CourseID:  strings.TrimSpace(courseID),
		StudentID: strings.TrimSpace(studentID),
	}

	if err := Validations.ValidateCourseStudent(courseStudent); err != nil {
		return nil, err
	}

	return courseStudent, nil
}

func ExistingCourseStudent(courseID, studentID string) (*Entities.CourseStudent, error) {
	return NewCourseStudent(courseID, studentID)
}
