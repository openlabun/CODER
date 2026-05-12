package cache_publisher

import (
	"context"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	ExamRepositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	cache_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
	PublisherAdapter "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/cache-handler/publisher"
)

type ExamRepository struct {
	publisherAdapter *PublisherAdapter.CachePublisherAdapter
	repository       ExamRepositories.ExamRepository
}

func NewExamRepository(publisherAdapter *PublisherAdapter.CachePublisherAdapter, repository ExamRepositories.ExamRepository) *ExamRepository {
	return &ExamRepository{publisherAdapter: publisherAdapter, repository: repository}
}

func (r *ExamRepository) CreateExam(ctx context.Context, exam *Entities.Exam) (*Entities.Exam, error) {
	message, err := PublisherAdapter.MapExamToSyncMessage(exam, cache_infrastructure.InsertOperation)
	if err != nil {
		return nil, err
	}
	if err := r.publisherAdapter.Save(message); err != nil {
		return nil, err
	}
	return exam, nil
}

func (r *ExamRepository) UpdateExam(ctx context.Context, exam *Entities.Exam) (*Entities.Exam, error) {
	message, err := PublisherAdapter.MapExamToSyncMessage(exam, cache_infrastructure.UpdateOperation)
	if err != nil {
		return nil, err
	}
	if err := r.publisherAdapter.Save(message); err != nil {
		return nil, err
	}
	return exam, nil
}

func (r *ExamRepository) DeleteExam(ctx context.Context, examID string) error {
	return r.repository.DeleteExam(ctx, examID)
}

func (r *ExamRepository) GetExamByID(ctx context.Context, examID string) (*Entities.Exam, error) {
	return r.repository.GetExamByID(ctx, examID)
}

func (r *ExamRepository) GetPublicExams(ctx context.Context, visibility string) ([]*Entities.Exam, error) {
	return r.repository.GetPublicExams(ctx, visibility)
}

func (r *ExamRepository) GetExamsByCourseID(ctx context.Context, courseID string) ([]*Entities.Exam, error) {
	return r.repository.GetExamsByCourseID(ctx, courseID)
}

func (r *ExamRepository) GetExamsByTeacherID(ctx context.Context, teacherID string) ([]*Entities.Exam, error) {
	return r.repository.GetExamsByTeacherID(ctx, teacherID)
}
