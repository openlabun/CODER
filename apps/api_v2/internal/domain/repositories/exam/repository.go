package exam_repository

import (
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
)

type ExamRepository interface {
	CreateExam(exam *Entities.Exam) (*Entities.Exam, error)
	UpdateExam(exam *Entities.Exam) (*Entities.Exam, error)
	DeleteExam(examID string) error

	GetExamByID(examID string) (*Entities.Exam, error)
	GetExamsByCourseID(courseID string) ([]*Entities.Exam, error)
	GetExamsByTeacherID(teacherID string) ([]*Entities.Exam, error)
}

type ChallengeRepository interface {
	CreateChallenge(challenge *Entities.Challenge) (*Entities.Challenge, error)
	UpdateChallenge(challenge *Entities.Challenge) (*Entities.Challenge, error)
	DeleteChallenge(challengeID string) error

	GetChallengeByID(challengeID string) (*Entities.Challenge, error)
	GetChallengesByExamID(examID string) ([]*Entities.Challenge, error)
	GetChallengesByTag(tag string) ([]*Entities.Challenge, error)
}

type TestCaseRepository interface {
	CreateTestCase(testCase *Entities.TestCase) (*Entities.TestCase, error)
	UpdateTestCase(testCase *Entities.TestCase) (*Entities.TestCase, error)
	DeleteTestCase(testCaseID string) error

	GetTestCaseByID(testCaseID string) (*Entities.TestCase, error)
	GetTestCasesByChallengeID(challengeID string) ([]*Entities.TestCase, error)
}
