package exam_validations

import (
	"fmt"
	"strings"
	"time"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
)

func ValidateExam(exam *Entities.Exam) error {
	if exam == nil {
		return fmt.Errorf("exam is nil")
	}

	if strings.TrimSpace(exam.ID) == "" {
		return fmt.Errorf("exam id is required")
	}
	if strings.TrimSpace(exam.Title) == "" {
		return fmt.Errorf("exam title is required")
	}
	if strings.TrimSpace(exam.ProfessorID) == "" {
		return fmt.Errorf("exam professor id is required")
	}

	if exam.Visibility == Entities.VisibilityCourse && (exam.CourseID == nil || strings.TrimSpace(*exam.CourseID) == "") {
		return fmt.Errorf("course visibility requires course id")
	}

	if exam.EndTime != nil && exam.EndTime.Before(exam.StartTime) {
		return fmt.Errorf("exam end time cannot be before start time")
	}

	if exam.TimeLimit < 0 {
		return fmt.Errorf("exam time limit cannot be negative")
	}
	if exam.TryLimit < 0 {
		return fmt.Errorf("exam try limit cannot be negative")
	}

	return nil
}

func ValidateExamEndTime (exam *Entities.Exam, now time.Time) error {
	if exam.EndTime != nil && now.After(*exam.EndTime) {
		return fmt.Errorf("exam can't end in the past")
	}

	if exam.EndTime != nil && exam.StartTime.After(*exam.EndTime) {
		return fmt.Errorf("exam can't end before it starts")
	}

	return nil
}

func ValidateExamTimeWindow(exam *Entities.Exam, now time.Time) error {
	if err := ValidateExam(exam); err != nil {
		return err
	}

	if now.Before(exam.StartTime) {
		return fmt.Errorf("exam has not started yet")
	}

	if exam.EndTime != nil && now.After(*exam.EndTime) && !exam.AllowLateSubmissions {
		return fmt.Errorf("exam is closed for submissions")
	}

	return nil
}
