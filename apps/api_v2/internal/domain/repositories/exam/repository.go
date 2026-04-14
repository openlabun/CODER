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
	GetPublicExams(ctx context.Context, visibility string) ([]*Entities.Exam, error)
	GetExamsByCourseID(ctx context.Context, courseID string) ([]*Entities.Exam, error)
	GetExamsByTeacherID(ctx context.Context, teacherID string) ([]*Entities.Exam, error)
}

type ExamScoreRepository interface {
	CreateExamScore(ctx context.Context, examScore *Entities.ExamScore) (*Entities.ExamScore, error)
	UpdateExamScore(ctx context.Context, examScore *Entities.ExamScore) (*Entities.ExamScore, error)
	DeleteExamScore(ctx context.Context, examScoreID string) error

	GetExamScores (ctx context.Context, examID, studentID *string) ([]*Entities.ExamScore, error)
	GetExamScoreByID(ctx context.Context, examScoreID string) (*Entities.ExamScore, error)
	GetExamScoresBySessionID(ctx context.Context, sessionID string) ([]*Entities.ExamScore, error)
}

type ChallengeRepository interface {
	CreateChallenge(ctx context.Context, challenge *Entities.Challenge) (*Entities.Challenge, error)
	UpdateChallenge(ctx context.Context, challenge *Entities.Challenge) (*Entities.Challenge, error)
	DeleteChallenge(ctx context.Context, challengeID string) error

	GetChallenges(ctx context.Context, status, tag, difficulty *string) ([]*Entities.Challenge, error)
	GetChallengeByID(ctx context.Context, challengeID string) (*Entities.Challenge, error)
	GetChallengesByExamID(ctx context.Context, examID string) ([]*Entities.Challenge, error)
	GetChallengesByUserID(ctx context.Context, userID string, examID *string) ([]*Entities.Challenge, error)
	GetChallengesByTag(ctx context.Context, tag string) ([]*Entities.Challenge, error)
	GetInputVariablesByChallengeID(ctx context.Context, challengeID string) ([]*Entities.IOVariable, error)
	GetOutputVariablesByChallengeID(ctx context.Context, challengeID string) ([]*Entities.IOVariable, error)
}

type ExamItemRepository interface {
	CreateExamItem(ctx context.Context, examItem *Entities.ExamItem) (*Entities.ExamItem, error)
	UpdateExamItem(ctx context.Context, examItem *Entities.ExamItem) (*Entities.ExamItem, error)
	DeleteExamItem(ctx context.Context, examItemID string) error

	GetExamItemByID(ctx context.Context, examItemID string) (*Entities.ExamItem, error)
	GetExamItem (ctx context.Context, examID *string, challengeID *string) ([]*Entities.ExamItem, error)
}

type ExamItemScoreRepository interface {
	CreateExamItemScore(ctx context.Context, examItemScore *Entities.ExamItemScore) (*Entities.ExamItemScore, error)
	UpdateExamItemScore(ctx context.Context, examItemScore *Entities.ExamItemScore) (*Entities.ExamItemScore, error)
	DeleteExamItemScore(ctx context.Context, examItemScoreID string) error

	GetExamItemScoreByID(ctx context.Context, examItemScoreID string) (*Entities.ExamItemScore, error)
	GetExamItemScoresByExamScoreID(ctx context.Context, examScoreID string) ([]*Entities.ExamItemScore, error)
	GetExamItemScore(ctx context.Context, examItemID, examScoreID string) (*Entities.ExamItemScore, error)
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

type IOVariableRepository interface {
	CreateIOVariable(ctx context.Context, ioVariable *Entities.IOVariable) (*Entities.IOVariable, error)
	UpdateIOVariable(ctx context.Context, ioVariable *Entities.IOVariable) (*Entities.IOVariable, error)
	GetIOVariableByID(ctx context.Context, ioVariableID string) (*Entities.IOVariable, error)
	DeleteIOVariable(ctx context.Context, ioVariableID string) error
}