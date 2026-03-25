package exam_repository

import (
	"context"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
)

type ExamRepository interface {
	CreateExam(ctx context.Context, exam *Entities.Exam) (*Entities.Exam, error)
	UpdateExam(ctx context.Context, exam *Entities.Exam) (*Entities.Exam, error)
	DeleteExam(ctx context.Context, examID string) error

	GetExamByID(ctx context.Context, examID string) (*Entities.Exam, error)
	GetExamsByCourseID(ctx context.Context, courseID string) ([]*Entities.Exam, error)
	GetExamsByTeacherID(ctx context.Context, teacherID string) ([]*Entities.Exam, error)
}

type ChallengeRepository interface {
	CreateChallenge(ctx context.Context, challenge *Entities.Challenge) (*Entities.Challenge, error)
	UpdateChallenge(ctx context.Context, challenge *Entities.Challenge) (*Entities.Challenge, error)
	DeleteChallenge(ctx context.Context, challengeID string) error

	GetChallengeByID(ctx context.Context, challengeID string) (*Entities.Challenge, error)
	GetChallengesByExamID(ctx context.Context, examID string) ([]*Entities.Challenge, error)
	GetChallengesByTag(ctx context.Context, tag string) ([]*Entities.Challenge, error)
	GetInputVariablesByChallengeID(ctx context.Context, challengeID string) ([]*Entities.IOVariable, error)
	GetOutputVariablesByChallengeID(ctx context.Context, challengeID string) ([]*Entities.IOVariable, error)
}

type ExamItemRepository interface {
	CreateExamItem(ctx context.Context, examItem *Entities.ExamItem) (*Entities.ExamItem, error)
	UpdateExamItem(ctx context.Context, examItem *Entities.ExamItem) (*Entities.ExamItem, error)
	GetExamItemByID(ctx context.Context, examItemID string) (*Entities.ExamItem, error)
	GetExamItem (ctx context.Context, examID *string, challengeID *string) ([]*Entities.ExamItem, error)
	DeleteExamItem(ctx context.Context, examItemID string) error
}

type TestCaseRepository interface {
	CreateTestCase(ctx context.Context, testCase *Entities.TestCase) (*Entities.TestCase, error)
	UpdateTestCase(ctx context.Context, testCase *Entities.TestCase) (*Entities.TestCase, error)
	DeleteTestCase(ctx context.Context, testCaseID string) error

	GetTestCaseByID(ctx context.Context, testCaseID string) (*Entities.TestCase, error)
	GetTestCasesByChallengeID(ctx context.Context, challengeID string) ([]*Entities.TestCase, error)
	GetInputVariablesByTestCaseID(ctx context.Context, testCaseID string) ([]*Entities.IOVariable, error)
	GetOutputVariablesByTestCaseID(ctx context.Context, testCaseID string) ([]*Entities.IOVariable, error)
}
