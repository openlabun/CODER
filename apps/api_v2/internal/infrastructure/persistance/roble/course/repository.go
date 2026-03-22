package roble_infrastructure

import (
	"context"
	"fmt"
	"strings"
	"time"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"
	infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
)

const (
	courseTableName        = "Course"
	courseStudentTableName = "CourseStudent"
)

type CourseRepository struct {
	adapter *infrastructure.RobleDatabaseAdapter
}

func NewCourseRepository(adapter *infrastructure.RobleDatabaseAdapter) *CourseRepository {
	return &CourseRepository{adapter: adapter}
}

func (r *CourseRepository) CreateCourse(ctx context.Context, course *Entities.Course) (*Entities.Course, error) {
	if course == nil {
		return nil, fmt.Errorf("course is nil")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	_, err := r.adapter.Insert(courseTableName, []map[string]any{courseToRecord(course)})
	if err != nil {
		return nil, err
	}

	return course, nil
}

func (r *CourseRepository) UpdateCourse(ctx context.Context, course *Entities.Course) (*Entities.Course, error) {
	if course == nil {
		return nil, fmt.Errorf("course is nil")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	courseID := course.ID

	updates := courseToUpdates(course)
	updates["UpdatedAt"] = time.Now().UTC().Format(time.RFC3339)

	_, err := r.adapter.Update(courseTableName, "ID", strings.TrimSpace(courseID), updates)
	if err != nil {
		return nil, err
	}

	course.ID = strings.TrimSpace(courseID)
	if ts, ok := updates["UpdatedAt"].(string); ok {
		if parsed, err := time.Parse(time.RFC3339, ts); err == nil {
			course.UpdatedAt = parsed
		}
	}

	return course, nil
}

func (r *CourseRepository) DeleteCourse(ctx context.Context, courseID string) error {
	if strings.TrimSpace(courseID) == "" {
		return fmt.Errorf("courseID is required")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return err
	}

	_, err := r.adapter.Delete(courseTableName, "ID", strings.TrimSpace(courseID))
	return err
}

func (r *CourseRepository) AddStudentToCourse(ctx context.Context, courseID, studentID string) error {
	if strings.TrimSpace(courseID) == "" {
		return fmt.Errorf("courseID is required")
	}
	if strings.TrimSpace(studentID) == "" {
		return fmt.Errorf("studentID is required")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return err
	}

	record := map[string]any{
		"CourseID":  strings.TrimSpace(courseID),
		"StudentID": strings.TrimSpace(studentID),
	}

	_, err := r.adapter.Insert(courseStudentTableName, []map[string]any{record})
	return err
}

func (r *CourseRepository) GetCourseByID(ctx context.Context, courseID string) (*Entities.Course, error) {
	if strings.TrimSpace(courseID) == "" {
		return nil, fmt.Errorf("courseID is required")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	res, err := r.adapter.Read(courseTableName, map[string]string{"ID": strings.TrimSpace(courseID)})
	if err != nil {
		return nil, err
	}

	record, err := firstRecord(res)
	if err != nil {
		return nil, nil
	}

	return recordToCourse(record)
}

func (r *CourseRepository) GetCourseByEnrollmentCode(ctx context.Context, enrollmentCode string) (*Entities.Course, error) {
	if strings.TrimSpace(enrollmentCode) == "" {
		return nil, fmt.Errorf("enrollmentCode is required")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	res, err := r.adapter.Read(courseTableName, map[string]string{"EnrollmentCode": strings.TrimSpace(enrollmentCode)})
	if err != nil {
		return nil, err
	}

	record, err := firstRecord(res)
	if err != nil {
		return nil, nil
	}

	return recordToCourse(record)
}

func (r *CourseRepository) GetCoursesByStudentID(ctx context.Context, studentID string) ([]*Entities.Course, error) {
	if strings.TrimSpace(studentID) == "" {
		return nil, fmt.Errorf("studentID is required")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	res, err := r.adapter.Read(courseStudentTableName, map[string]string{"StudentID": strings.TrimSpace(studentID)})
	if err != nil {
		return nil, err
	}

	relations := extractRecords(res)
	if len(relations) == 0 {
		return []*Entities.Course{}, nil
	}

	courses := make([]*Entities.Course, 0, len(relations))
	for _, relation := range relations {
		courseID := asString(relation["CourseID"])
		if courseID == "" {
			continue
		}

		course, err := r.GetCourseByID(ctx, courseID)
		if err != nil {
			return nil, err
		}
		if course != nil {
			courses = append(courses, course)
		}
	}

	return courses, nil
}

func (r *CourseRepository) GetCoursesByTeacherID(ctx context.Context, teacherID string) ([]*Entities.Course, error) {
	if strings.TrimSpace(teacherID) == "" {
		return nil, fmt.Errorf("teacherID is required")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	res, err := r.adapter.Read(courseTableName, map[string]string{"ProfessorID": strings.TrimSpace(teacherID)})
	if err != nil {
		return nil, err
	}

	records := extractRecords(res)
	if len(records) == 0 {
		return []*Entities.Course{}, nil
	}

	courses := make([]*Entities.Course, 0, len(records))
	for _, record := range records {
		course, err := recordToCourse(record)
		if err != nil {
			return nil, err
		}
		if course != nil {
			courses = append(courses, course)
		}
	}

	return courses, nil
}


func firstRecord(res map[string]any) (map[string]any, error) {
	if data, ok := res["data"]; ok {
		if arr := asRecordSlice(data); len(arr) > 0 {
			return arr[0], nil
		}
	}

	if records, ok := res["records"]; ok {
		if arr := asRecordSlice(records); len(arr) > 0 {
			return arr[0], nil
		}
	}

	if arr := asRecordSlice(res); len(arr) > 0 {
		return arr[0], nil
	}

	if isCourseRecordShape(res) {
		return res, nil
	}

	return nil, fmt.Errorf("no records found")
}

func extractRecords(res map[string]any) []map[string]any {
	if data, ok := res["data"]; ok {
		if arr := asRecordSlice(data); len(arr) > 0 {
			return arr
		}
	}

	if records, ok := res["records"]; ok {
		if arr := asRecordSlice(records); len(arr) > 0 {
			return arr
		}
	}

	if arr := asRecordSlice(res); len(arr) > 0 {
		return arr
	}

	if isCourseRecordShape(res) {
		return []map[string]any{res}
	}

	return nil
}

func asRecordSlice(value any) []map[string]any {
	items, ok := value.([]any)
	if !ok {
		return nil
	}

	results := make([]map[string]any, 0, len(items))
	for _, item := range items {
		if row, ok := item.(map[string]any); ok {
			results = append(results, row)
		}
	}

	return results
}

func isCourseRecordShape(record map[string]any) bool {
	_, hasID := record["ID"]
	_, hasName := record["Name"]
	return hasID && hasName
}

func asString(v any) string {
	if v == nil {
		return ""
	}

	if s, ok := v.(string); ok {
		return strings.TrimSpace(s)
	}

	return strings.TrimSpace(fmt.Sprintf("%v", v))
}

func asInt(v any) int {
	switch value := v.(type) {
	case int:
		return value
	case int32:
		return int(value)
	case int64:
		return int(value)
	case float64:
		return int(value)
	case string:
		var parsed int
		_, err := fmt.Sscanf(strings.TrimSpace(value), "%d", &parsed)
		if err == nil {
			return parsed
		}
	}

	return 0
}

func asTime(v any) (time.Time, bool) {
	s := asString(v)
	if s == "" {
		return time.Time{}, false
	}

	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, true
	}
	if t, err := time.Parse("2006-01-02T15:04", s); err == nil {
		return t, true
	}
	if t, err := time.Parse("2006-01-02 15:04:05", s); err == nil {
		return t, true
	}

	return time.Time{}, false
}
