package submission_repository

import (
	"context"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
)

type SubmissionRepository interface {
	CreateSubmission(ctx context.Context, submission *Entities.Submission) (*Entities.Submission, error)
	UpdateSubmission(ctx context.Context, submission *Entities.Submission) (*Entities.Submission, error)
	DeleteSubmission(ctx context.Context, submissionID string) error

	GetSubmissionByID(ctx context.Context, submissionID string) (*Entities.Submission, error)
	GetSubmissionsByExamItemScoreID(ctx context.Context, examItemScoreID string) ([]*Entities.Submission, error)
	GetBestSubmissionByExamItemScoreID(ctx context.Context, examItemScoreID string) (*Entities.Submission, error)
	GetLastSubmissionByExamItemScoreID(ctx context.Context, examItemScoreID string) (*Entities.Submission, error)
	GetSubmissionsBySessionID(ctx context.Context, sessionID string, status *string, testID *string, challengeID *string) ([]*Entities.Submission, error)
	GetSubmissionsByUserID(ctx context.Context, userID string, status *string, testID *string, challengeID *string) ([]*Entities.Submission, error)
	GetSubmissionsByChallengeID(ctx context.Context, challengeID string, status *string, testID *string) ([]*Entities.Submission, error)
}

type SessionRepository interface {
	CreateSession(ctx context.Context, session *Entities.Session) (*Entities.Session, error)
	UpdateSession(ctx context.Context, session *Entities.Session) (*Entities.Session, error)
	DeleteSession(ctx context.Context, sessionID string) error

	GetSessionByID(ctx context.Context, sessionID string) (*Entities.Session, error)
	GetSessionsByExamID(ctx context.Context, examID string) ([]*Entities.Session, error)
	GetSessionsByStudentID(ctx context.Context, studentID string) ([]*Entities.Session, error)
}

type SubmissionResultRepository interface {
	CreateResult(ctx context.Context, result *Entities.SubmissionResult) (*Entities.SubmissionResult, error)
	UpdateResult(ctx context.Context, result *Entities.SubmissionResult) (*Entities.SubmissionResult, error)
	DeleteResult(ctx context.Context, resultID string) error

	GetResultByID(ctx context.Context, resultID string) (*Entities.SubmissionResult, error)
	GetResultsBySubmissionID(ctx context.Context, submissionID string) ([]*Entities.SubmissionResult, error)
	GetResultByTestCase(ctx context.Context, testCaseID string) ([]*Entities.SubmissionResult, error)
}
