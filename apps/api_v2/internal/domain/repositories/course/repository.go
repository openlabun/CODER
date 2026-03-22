package course_repository

import (
	"context"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"
)

type CourseRepository interface {
	CreateCourse(ctx context.Context, course *Entities.Course) (*Entities.Course, error)
	UpdateCourse(ctx context.Context, course *Entities.Course) (*Entities.Course, error)
	DeleteCourse(ctx context.Context, courseID string) error
	AddStudentToCourse(ctx context.Context, courseID, studentID string) error

	GetCourseByID(ctx context.Context, courseID string) (*Entities.Course, error)
	GetCourseByEnrollmentCode(ctx context.Context, enrollmentCode string) (*Entities.Course, error)
	GetCoursesByStudentID(ctx context.Context, studentID string) ([]*Entities.Course, error)
	GetCoursesByTeacherID(ctx context.Context, teacherID string) ([]*Entities.Course, error)
}
