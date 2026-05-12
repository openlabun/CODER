package cache_roble_infrastructure

import (
	"context"
	"strings"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	cache "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
	roble_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
	repository "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/exam"
)

type ExamRepository struct {
	robleRepository *repository.ExamRepository
	cacheAdapter    *cache.CacheAdapter
}

func NewExamRepository(robleAdapter *roble_infrastructure.RobleDatabaseAdapter, cacheAdapter *cache.CacheAdapter) *ExamRepository {
	return &ExamRepository{
		robleRepository: repository.NewExamRepository(robleAdapter),
		cacheAdapter:    cacheAdapter,
	}
}

func (r *ExamRepository) CreateExam(ctx context.Context, exam *Entities.Exam) (*Entities.Exam, error) {
	result, err := r.robleRepository.CreateExam(ctx, exam)
	if err != nil {
		return nil, err
	}
	if rec, e := cache.EntityToMap(result); e == nil {
		r.cacheAdapter.Save(cache.CacheExamTable, cache.InsertOperation, rec)
	}
	return result, nil
}

func (r *ExamRepository) UpdateExam(ctx context.Context, exam *Entities.Exam) (*Entities.Exam, error) {
	result, err := r.robleRepository.UpdateExam(ctx, exam)
	if err != nil {
		return nil, err
	}
	if rec, e := cache.EntityToMap(result); e == nil {
		r.cacheAdapter.Save(cache.CacheExamTable, cache.UpdateOperation, rec)
	}
	return result, nil
}

func (r *ExamRepository) DeleteExam(ctx context.Context, examID string) error {
	if err := r.robleRepository.DeleteExam(ctx, examID); err != nil {
		return err
	}
	r.cacheAdapter.DeleteByID(cache.CacheExamTable, strings.TrimSpace(examID))
	return nil
}

func (r *ExamRepository) GetExamByID(ctx context.Context, examID string) (*Entities.Exam, error) {
	id := strings.TrimSpace(examID)

	if rec, err := r.cacheAdapter.FindByID(cache.CacheExamTable, id); err == nil && rec != nil {
		cache.NormalizeBools(rec, "allow_late_submissions")
		if exam, e := cache.MapToEntity[Entities.Exam](rec); e == nil {
			return exam, nil
		}
	}

	exam, err := r.robleRepository.GetExamByID(ctx, examID)
	if err != nil || exam == nil {
		return exam, err
	}
	if rec, e := cache.EntityToMap(exam); e == nil {
		r.cacheAdapter.Save(cache.CacheExamTable, cache.InsertOperation, rec)
	}
	return exam, nil
}

func (r *ExamRepository) GetPublicExams(ctx context.Context, visibility string) ([]*Entities.Exam, error) {
	exams, err := r.robleRepository.GetPublicExams(ctx, visibility)
	if err != nil {
		return nil, err
	}
	for _, exam := range exams {
		if rec, e := cache.EntityToMap(exam); e == nil {
			r.cacheAdapter.Save(cache.CacheExamTable, cache.InsertOperation, rec)
		}
	}
	return exams, nil
}

func (r *ExamRepository) GetExamsByCourseID(ctx context.Context, courseID string) ([]*Entities.Exam, error) {
	id := strings.TrimSpace(courseID)
	if recs, err := r.cacheAdapter.FindBy(cache.CacheExamTable, map[string]string{"course_id": id}); err == nil && recs != nil {
		return examsFromRecords(recs), nil
	}
	exams, err := r.robleRepository.GetExamsByCourseID(ctx, courseID)
	if err != nil {
		return nil, err
	}
	for _, exam := range exams {
		if rec, e := cache.EntityToMap(exam); e == nil {
			r.cacheAdapter.Save(cache.CacheExamTable, cache.InsertOperation, rec)
		}
	}
	return exams, nil
}

func (r *ExamRepository) GetExamsByTeacherID(ctx context.Context, teacherID string) ([]*Entities.Exam, error) {
	id := strings.TrimSpace(teacherID)
	if recs, err := r.cacheAdapter.FindBy(cache.CacheExamTable, map[string]string{"professor_id": id}); err == nil && recs != nil {
		return examsFromRecords(recs), nil
	}
	exams, err := r.robleRepository.GetExamsByTeacherID(ctx, teacherID)
	if err != nil {
		return nil, err
	}
	for _, exam := range exams {
		if rec, e := cache.EntityToMap(exam); e == nil {
			r.cacheAdapter.Save(cache.CacheExamTable, cache.InsertOperation, rec)
		}
	}
	return exams, nil
}