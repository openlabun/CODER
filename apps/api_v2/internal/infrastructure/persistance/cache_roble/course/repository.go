package cache_roble_infrastructure

import (
	"context"
	"strings"

	UserEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"

	cache "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
	roble_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
	repository "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/course"
)

type CourseRepository struct {
	robleRepository *repository.CourseRepository
	cacheAdapter    *cache.CacheAdapter
}

func NewCourseRepository(robleAdapter *roble_infrastructure.RobleDatabaseAdapter, cacheAdapter *cache.CacheAdapter) *CourseRepository {
	return &CourseRepository{
		robleRepository: repository.NewCourseRepository(robleAdapter),
		cacheAdapter:    cacheAdapter,
	}
}

func (r *CourseRepository) CreateCourse(ctx context.Context, course *Entities.Course) (*Entities.Course, error) {
	result, err := r.robleRepository.CreateCourse(ctx, course)
	if err != nil {
		return nil, err
	}
	if rec, e := cache.EntityToMap(result); e == nil {
		r.cacheAdapter.Save(cache.CacheCourseTable, cache.InsertOperation, rec)
	}
	return result, nil
}

func (r *CourseRepository) UpdateCourse(ctx context.Context, course *Entities.Course) (*Entities.Course, error) {
	result, err := r.robleRepository.UpdateCourse(ctx, course)
	if err != nil {
		return nil, err
	}
	if rec, e := cache.EntityToMap(result); e == nil {
		r.cacheAdapter.Save(cache.CacheCourseTable, cache.UpdateOperation, rec)
	}
	return result, nil
}

func (r *CourseRepository) DeleteCourse(ctx context.Context, courseID string) error {
	if err := r.robleRepository.DeleteCourse(ctx, courseID); err != nil {
		return err
	}
	r.cacheAdapter.DeleteByID(cache.CacheCourseTable, strings.TrimSpace(courseID))
	return nil
}

func (r *CourseRepository) AddStudentToCourse(ctx context.Context, courseID string, studentID string) error {
	return r.robleRepository.AddStudentToCourse(ctx, courseID, studentID)
}

func (r *CourseRepository) RemoveStudentFromCourse(ctx context.Context, courseID string, studentID string) error {
	return r.robleRepository.RemoveStudentFromCourse(ctx, courseID, studentID)
}

func (r *CourseRepository) GetCourseByID(ctx context.Context, courseID string) (*Entities.Course, error) {
	id := strings.TrimSpace(courseID)
	if rec, err := r.cacheAdapter.FindByID(cache.CacheCourseTable, id); err == nil && rec != nil {
		if c, e := cache.MapToEntity[Entities.Course](rec); e == nil {
			return c, nil
		}
	}
	c, err := r.robleRepository.GetCourseByID(ctx, courseID)
	if err != nil || c == nil {
		return c, err
	}
	if rec, e := cache.EntityToMap(c); e == nil {
		r.cacheAdapter.Save(cache.CacheCourseTable, cache.InsertOperation, rec)
	}
	return c, nil
}

func (r *CourseRepository) GetCoursesByEnrollmentCode(ctx context.Context, enrollmentCode string) (*Entities.Course, error) {
	code := strings.TrimSpace(enrollmentCode)
	if recs, err := r.cacheAdapter.FindBy(cache.CacheCourseTable, map[string]string{"enrollment_code": code}); err == nil && len(recs) > 0 {
		if c, e := cache.MapToEntity[Entities.Course](recs[0]); e == nil {
			return c, nil
		}
	}
	c, err := r.robleRepository.GetCourseByEnrollmentCode(ctx, enrollmentCode)
	if err != nil || c == nil {
		return c, err
	}
	if rec, e := cache.EntityToMap(c); e == nil {
		r.cacheAdapter.Save(cache.CacheCourseTable, cache.InsertOperation, rec)
	}
	return c, nil
}

func (r *CourseRepository) GetCoursesByStudentID(ctx context.Context, studentID string) ([]*Entities.Course, error) {
	id := strings.TrimSpace(studentID)
	if recs, err := r.cacheAdapter.FindBy(cache.CacheCourseTable, map[string]string{"student_id": id}); err == nil && recs != nil {
		return coursesFromRecords(recs), nil
	}
	courses, err := r.robleRepository.GetCoursesByStudentID(ctx, studentID)
	if err != nil {
		return nil, err
	}
	for _, c := range courses {
		if rec, e := cache.EntityToMap(c); e == nil {
			r.cacheAdapter.Save(cache.CacheCourseTable, cache.InsertOperation, rec)
		}
	}
	return courses, nil
}

func (r *CourseRepository) GetCoursesByTeacherID(ctx context.Context, teacherID string) ([]*Entities.Course, error) {
	id := strings.TrimSpace(teacherID)
	if recs, err := r.cacheAdapter.FindBy(cache.CacheCourseTable, map[string]string{"teacher_id": id}); err == nil && recs != nil {
		return coursesFromRecords(recs), nil
	}
	courses, err := r.robleRepository.GetCoursesByTeacherID(ctx, teacherID)
	if err != nil {
		return nil, err
	}
	for _, c := range courses {
		if rec, e := cache.EntityToMap(c); e == nil {
			r.cacheAdapter.Save(cache.CacheCourseTable, cache.InsertOperation, rec)
		}
	}
	return courses, nil
}

func (r *CourseRepository) GetStudentsByCourseID(ctx context.Context, courseID string) ([]*UserEntities.User, error) {
	return r.robleRepository.GetStudentsByCourseID(ctx, courseID)
}