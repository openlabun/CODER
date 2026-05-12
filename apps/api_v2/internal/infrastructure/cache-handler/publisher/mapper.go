package cache_publisher

import (
	"encoding/json"

	CourseEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"
	ExamEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	SubmissionEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	UserEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"

	cache_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
)

func mapToSyncMessage(model string, requestType string, entity any) (SyncMessage, error) {
	data, err := json.Marshal(entity)
	if err != nil {
		return SyncMessage{}, err
	}
	return SyncMessage{Model: model, Request: requestType, Data: data}, nil
}

func MapUserToSyncMessage(user *UserEntities.User, requestType string) (SyncMessage, error) {
	return mapToSyncMessage(cache_infrastructure.UserEntity, requestType, user)
}

func MapCourseToSyncMessage(course *CourseEntities.Course, requestType string) (SyncMessage, error) {
	return mapToSyncMessage(cache_infrastructure.CourseEntity, requestType, course)
}

func MapExamToSyncMessage(exam *ExamEntities.Exam, requestType string) (SyncMessage, error) {
	return mapToSyncMessage(cache_infrastructure.ExamEntity, requestType, exam)
}

func MapChallengeToSyncMessage(challenge *ExamEntities.Challenge, requestType string) (SyncMessage, error) {
	return mapToSyncMessage(cache_infrastructure.ChallengeEntity, requestType, challenge)
}

func MapExamItemToSyncMessage(examItem *ExamEntities.ExamItem, requestType string) (SyncMessage, error) {
	return mapToSyncMessage(cache_infrastructure.ExamItemEntity, requestType, examItem)
}

func MapExamItemScoreToSyncMessage(examItemScore *ExamEntities.ExamItemScore, requestType string) (SyncMessage, error) {
	return mapToSyncMessage(cache_infrastructure.ExamItemScoreEntity, requestType, examItemScore)
}

func MapExamScoreToSyncMessage(examScore *ExamEntities.ExamScore, requestType string) (SyncMessage, error) {
	return mapToSyncMessage(cache_infrastructure.ExamScoreEntity, requestType, examScore)
}

func MapIOVariableToSyncMessage(ioVariable *ExamEntities.IOVariable, requestType string) (SyncMessage, error) {
	return mapToSyncMessage(cache_infrastructure.IOVariableEntity, requestType, ioVariable)
}

func MapTestCaseToSyncMessage(testCase *ExamEntities.TestCase, requestType string) (SyncMessage, error) {
	return mapToSyncMessage(cache_infrastructure.TestCaseEntity, requestType, testCase)
}

func MapSessionToSyncMessage(session *SubmissionEntities.Session, requestType string) (SyncMessage, error) {
	return mapToSyncMessage(cache_infrastructure.SessionEntity, requestType, session)
}

func MapSubmissionToSyncMessage(submission *SubmissionEntities.Submission, requestType string) (SyncMessage, error) {
	return mapToSyncMessage(cache_infrastructure.SubmissionEntity, requestType, submission)
}

func MapSubmissionResultToSyncMessage(submissionResult *SubmissionEntities.SubmissionResult, requestType string) (SyncMessage, error) {
	return mapToSyncMessage(cache_infrastructure.SubmissionResultEntity, requestType, submissionResult)
}
