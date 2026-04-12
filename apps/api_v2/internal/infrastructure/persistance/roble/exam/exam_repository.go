package roble_infrastructure

import (
	"context"
	"fmt"
	"strings"
	"time"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
)

const (
	examTableName = "Exam"
)

type ExamRepository struct {
	adapter *infrastructure.RobleDatabaseAdapter
}

func NewExamRepository(adapter *infrastructure.RobleDatabaseAdapter) *ExamRepository {
	return &ExamRepository{adapter: adapter}
}

func (r *ExamRepository) CreateExam(ctx context.Context, exam *Entities.Exam) (*Entities.Exam, error) {
	if exam == nil {
		return nil, fmt.Errorf("exam is nil")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	_, err := r.adapter.Insert(examTableName, []map[string]any{examToRecord(exam)})
	if err != nil {
		return nil, err
	}

	return exam, nil
}

func (r *ExamRepository) UpdateExam(ctx context.Context, exam *Entities.Exam) (*Entities.Exam, error) {
	if exam == nil {
		return nil, fmt.Errorf("exam is nil")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	examID := strings.TrimSpace(exam.ID)
	if examID == "" {
		return nil, fmt.Errorf("exam id is required")
	}

	updates := examToUpdates(exam)
	updates["UpdatedAt"] = time.Now().UTC().Format(time.RFC3339)

	_, err := r.adapter.Update(examTableName, "ID", examID, updates)
	if err != nil {
		return nil, err
	}

	exam.ID = examID
	if ts, ok := updates["UpdatedAt"].(string); ok {
		if parsed, parseErr := time.Parse(time.RFC3339, ts); parseErr == nil {
			exam.UpdatedAt = parsed
		}
	}

	return exam, nil
}

func (r *ExamRepository) DeleteExam(ctx context.Context, examID string) error {
	normalizedID := strings.TrimSpace(examID)
	if normalizedID == "" {
		return fmt.Errorf("examID is required")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return err
	}

	_, err := r.adapter.Delete(examTableName, "ID", normalizedID)
	return err
}

func (r *ExamRepository) GetExamByID(ctx context.Context, examID string) (*Entities.Exam, error) {
	normalizedID := strings.TrimSpace(examID)
	if normalizedID == "" {
		return nil, fmt.Errorf("examID is required")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	res, err := r.adapter.Read(examTableName, map[string]string{"ID": normalizedID})
	if err != nil {
		return nil, err
	}

	record, err := firstRecord(res)
	if err != nil {
		return nil, nil
	}

	return recordToExam(record)
}

func (r *ExamRepository) GetPublicExams(ctx context.Context, visibility string) ([]*Entities.Exam, error) {
	normalized := strings.TrimSpace(visibility)
	if normalized == "" {
		return nil, fmt.Errorf("visibility is required")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	res, err := r.adapter.Read(examTableName, map[string]string{"Visibility": normalized})
	if err != nil {
		return nil, err
	}

	records := extractRecords(res)
	if len(records) == 0 {
		return []*Entities.Exam{}, nil
	}

	// TODO: Replace with null comparison in DB query
	public_records := make([]map[string]any, 0, len(records))
	for _, record := range records {
		if record["CourseID"] == nil || strings.TrimSpace(fmt.Sprint(record["CourseID"])) == "" {
			public_records = append(public_records, record)
		}
	}

	exams := make([]*Entities.Exam, 0, len(public_records))
	for _, record := range public_records {
		exam, mapErr := recordToExam(record)
		if mapErr != nil {
			// Ignore malformed records to avoid failing the whole listing.
			// This is specially relevant for legacy records with invalid visibility/course combinations.
			continue
		}
		if exam != nil {
			exams = append(exams, exam)
		}
	}

	return exams, nil
}

func (r *ExamRepository) GetExamsByCourseID(ctx context.Context, courseID string) ([]*Entities.Exam, error) {
	normalizedID := strings.TrimSpace(courseID)
	if normalizedID == "" {
		return nil, fmt.Errorf("courseID is required")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	res, err := r.adapter.Read(examTableName, map[string]string{"CourseID": normalizedID})
	if err != nil {
		return nil, err
	}

	records := extractRecords(res)
	if len(records) == 0 {
		return []*Entities.Exam{}, nil
	}

	exams := make([]*Entities.Exam, 0, len(records))
	for _, record := range records {
		exam, mapErr := recordToExam(record)
		if mapErr != nil {
			// Ignore malformed records to avoid failing the whole listing.
			// This is specially relevant for legacy records with invalid visibility/course combinations.
			continue
		}
		if exam != nil {
			exams = append(exams, exam)
		}
	}

	return exams, nil
}

func (r *ExamRepository) GetExamsByTeacherID(ctx context.Context, teacherID string) ([]*Entities.Exam, error) {
	normalizedID := strings.TrimSpace(teacherID)
	if normalizedID == "" {
		return nil, fmt.Errorf("teacherID is required")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	res, err := r.adapter.Read(examTableName, map[string]string{"ProfessorID": normalizedID})
	if err != nil {
		return nil, err
	}

	records := extractRecords(res)
	if len(records) == 0 {
		return []*Entities.Exam{}, nil
	}

	exams := make([]*Entities.Exam, 0, len(records))
	for _, record := range records {
		exam, mapErr := recordToExam(record)
		if mapErr != nil {
			// Ignore malformed records to avoid failing the whole listing.
			// This is specially relevant for legacy records with invalid visibility/course combinations.
			continue
		}
		if exam != nil {
			exams = append(exams, exam)
		}
	}

	return exams, nil
}
