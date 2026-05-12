package cache_publisher

import (
	"context"

	CourseRepositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/course"
	UserEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"
	
	cache_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
	PublisherAdapter "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/cache-handler/publisher"
)

type CourseRepository struct {
	publisherAdapter *PublisherAdapter.CachePublisherAdapter
	repository 	 	 CourseRepositories.CourseRepository
}

func NewCourseRepository(publisherAdapter *PublisherAdapter.CachePublisherAdapter, repository CourseRepositories.CourseRepository) *CourseRepository {
	return &CourseRepository{publisherAdapter: publisherAdapter, repository: repository}
}

func (r *CourseRepository) CreateCourse(ctx context.Context, course *Entities.Course) (*Entities.Course, error) {
	// [STEP 1] Create SyncMessage from the course entity
	message, err := PublisherAdapter.MapCourseToSyncMessage(course, cache_infrastructure.InsertOperation)
	if err != nil {
		return nil, err
	}

	// [STEP 2] Publish the SyncMessage to RabbitMQ
	if err := r.publisherAdapter.Save(message); err != nil {
		return nil, err
	}
	
	return course, nil
}

func (r *CourseRepository) UpdateCourse(ctx context.Context, course *Entities.Course) (*Entities.Course, error) {
	// [STEP 1] Create SyncMessage from the course entity
	message, err := PublisherAdapter.MapCourseToSyncMessage(course, cache_infrastructure.UpdateOperation)
	if err != nil {
		return nil, err
	}

	// [STEP 2] Publish the SyncMessage to RabbitMQ
	if err := r.publisherAdapter.Save(message); err != nil {
		return nil, err
	}

	return course, nil
}

func (r *CourseRepository) DeleteCourse(ctx context.Context, courseID string) error {
	// Wrapper
	return r.repository.DeleteCourse(ctx, courseID)
}

func (r *CourseRepository) AddStudentToCourse(ctx context.Context, courseID string, studentID string) error {
	// Wrapper
	return r.repository.AddStudentToCourse(ctx, courseID, studentID)
}

func (r *CourseRepository) RemoveStudentFromCourse(ctx context.Context, courseID string, studentID string) error {
	// Wrapper
	return r.repository.RemoveStudentFromCourse(ctx, courseID, studentID)
}

func (r *CourseRepository) GetCourseByID(ctx context.Context, courseID string) (*Entities.Course, error) {
	// Wrapper
	return r.repository.GetCourseByID(ctx, courseID)
}

func (r *CourseRepository) GetCoursesByEnrollmentCode(ctx context.Context, enrollmentCode string) (*Entities.Course, error) {
	// Wrapper
	return r.repository.GetCourseByEnrollmentCode(ctx, enrollmentCode)
}

func (r *CourseRepository) GetCoursesByStudentID(ctx context.Context, studentID string) ([]*Entities.Course, error) {
	// Wrapper
	return r.repository.GetCoursesByStudentID(ctx, studentID)
}

func (r *CourseRepository) GetCoursesByTeacherID(ctx context.Context, teacherID string) ([]*Entities.Course, error) {
	// Wrapper
	return r.repository.GetCoursesByTeacherID(ctx, teacherID)
}

func (r *CourseRepository) GetStudentsByCourseID(ctx context.Context, courseID string) ([]*UserEntities.User, error) {
	// Wrapper
	return r.repository.GetStudentsByCourseID(ctx, courseID)
}




