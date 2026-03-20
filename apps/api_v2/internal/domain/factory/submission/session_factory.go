package submission_factory

import (
	"strings"
	"time"

	Entities "../../entities/submission"
	ExamEntities "../../entities/exam"
	Validations "../../validations/submission"
)

func NewSession(id, studentID string, exam *ExamEntities.Exam) (*Entities.Session, error) {
	now := time.Now()

	timeLeft := 0
	if exam.TimeLimit < 0 {
		timeLeft = exam.TimeLimit
	}
	
	session := &Entities.Session{
		ID:            strings.TrimSpace(id),
		StudentID:     strings.TrimSpace(studentID),
		ExamID:        strings.TrimSpace(exam.ID),
		Status:        Entities.SessionStatusActive,
		Attempts:      0,
		TimeLeft:      timeLeft,
		StartedAt:     now,
		FinishedAt:    nil,
		LastHeartbeat: now,
	}

	if err := Validations.ValidateSession(session); err != nil {
		return nil, err
	}

	return session, nil
}

func ExistingSession(
	id, studentID, examID string,
	status Entities.SessionStatus,
	attempts, timeLeft int,
	startedAt time.Time,
	finishedAt *time.Time,
	lastHeartbeat time.Time,
) (*Entities.Session, error) {
	session := &Entities.Session{
		ID:            strings.TrimSpace(id),
		StudentID:     strings.TrimSpace(studentID),
		ExamID:        strings.TrimSpace(examID),
		Status:        status,
		Attempts:      attempts,
		TimeLeft:      timeLeft,
		StartedAt:     startedAt,
		FinishedAt:    finishedAt,
		LastHeartbeat: lastHeartbeat,
	}

	if err := Validations.ValidateSession(session); err != nil {
		return nil, err
	}

	return session, nil
}