package submission_repository

import (
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
)

type SubmissionRepository interface {
	CreateSubmission(submission *Entities.Submission) (*Entities.Submission, error)
	UpdateSubmission(submission *Entities.Submission) (*Entities.Submission, error)
	DeleteSubmission(submissionID string) error

	GetSubmissionByID(submissionID string) (*Entities.Submission, error)
	GetSubmissionsBySessionID(sessionID string) ([]*Entities.Submission, error)
	GetSubmissionsByUserID(userID string) ([]*Entities.Submission, error)
	GetSubmissionsByChallengeID(challengeID string) ([]*Entities.Submission, error)
}

type SessionRepository interface {
	CreateSession(session *Entities.Session) (*Entities.Session, error)
	UpdateSession(session *Entities.Session) (*Entities.Session, error)
	DeleteSession(sessionID string) error

	GetSessionByID(sessionID string) (*Entities.Session, error)
	GetSessionsByExamID(examID string) ([]*Entities.Session, error)
	GetSessionsByStudentID(studentID string) ([]*Entities.Session, error)
}

type SubmissionResultRepository interface {
	CreateResult(result *Entities.SubmissionResult) (*Entities.SubmissionResult, error)
	UpdateResult(result *Entities.SubmissionResult) (*Entities.SubmissionResult, error)
	DeleteResult(resultID string) error

	GetResultByID(resultID string) (*Entities.SubmissionResult, error)
	GetResultsBySubmissionID(submissionID string) ([]*Entities.SubmissionResult, error)
	GetResultByTestCase(testCaseID string) ([]*Entities.SubmissionResult, error)
}
