package course_repository

import (
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"
)

type CourseRepository interface {
	CreateCourse(course *Entities.Course) (*Entities.Course, error)
	UpdateCourse(course *Entities.Course) (*Entities.Course, error)
	DeleteCourse(courseID string) error
	AddStudentToCourse(courseID, studentID string) error

	GetCourseByID(courseID string) (*Entities.Course, error)
	GetCourseByEnrollmentCode(enrollmentCode string) (*Entities.Course, error)
	GetCoursesByStudentID(studentID string) ([]*Entities.Course, error)
	GetCoursesByTeacherID(teacherID string) ([]*Entities.Course, error)
}
