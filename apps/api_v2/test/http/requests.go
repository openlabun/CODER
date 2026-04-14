package http_tests

import (
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func PostAuthRegister(t *testing.T, app *fiber.App, email, name, password string) *HTTPResponse {
	body := map[string]any{"email": email, "name": name, "password": password}
	return doJSONRequest(t, app, "POST", "/auth/register", body, nil)
}

func PostAuthLogin(t *testing.T, app *fiber.App, email, password string) *HTTPResponse {
	body := map[string]any{"email": email, "password": password}
	return doJSONRequest(t, app, "POST", "/auth/login", body, nil)
}

func GetAuthMe(t *testing.T, app *fiber.App, access *HTTPAccess, userID string) *HTTPResponse {
	path := "/auth/me"
	if userID != "" {
		path = fmt.Sprintf("/auth/me?userId=%s", userID)
	}
	return doJSONRequest(t, app, "GET", path, nil, authHeaders(access))
}

func PostCourseCreate(t *testing.T, app *fiber.App, access *HTTPAccess, body map[string]any) *HTTPResponse {
	return doJSONRequest(t, app, "POST", "/courses/", body, authHeaders(access))
}

func PostCourseUpdate(t *testing.T, app *fiber.App, access *HTTPAccess, courseID string, body map[string]any) *HTTPResponse {
	path := fmt.Sprintf("/courses/%s", courseID)
	return doJSONRequest(t, app, "POST", path, body, authHeaders(access))
}

func GetCourseByID(t *testing.T, app *fiber.App, access *HTTPAccess, courseID string) *HTTPResponse {
	path := fmt.Sprintf("/courses/%s", courseID)
	return doJSONRequest(t, app, "GET", path, nil, authHeaders(access))
}

func DeleteCourseByID(t *testing.T, app *fiber.App, access *HTTPAccess, courseID string) *HTTPResponse {
	path := fmt.Sprintf("/courses/%s", courseID)
	return doJSONRequest(t, app, "DELETE", path, nil, authHeaders(access))
}

func PostCourseEnroll(t *testing.T, app *fiber.App, access *HTTPAccess, body map[string]any) *HTTPResponse {
	return doJSONRequest(t, app, "POST", "/courses/enroll", body, authHeaders(access))
}

func GetCourses(t *testing.T, app *fiber.App, access *HTTPAccess, scope, studentID, teacherID string) *HTTPResponse {
	path := fmt.Sprintf("/courses/?scope=%s", scope)
	if studentID != "" {
		path = fmt.Sprintf("%s&studentId=%s", path, studentID)
	}
	if teacherID != "" {
		path = fmt.Sprintf("%s&teacherId=%s", path, teacherID)
	}
	return doJSONRequest(t, app, "GET", path, nil, authHeaders(access))
}

func GetCourseStudents(t *testing.T, app *fiber.App, access *HTTPAccess, courseID string) *HTTPResponse {
	path := fmt.Sprintf("/courses/%s/students", courseID)
	return doJSONRequest(t, app, "GET", path, nil, authHeaders(access))
}

func DeleteCourseStudent(t *testing.T, app *fiber.App, access *HTTPAccess, courseID, studentID string) *HTTPResponse {
	path := fmt.Sprintf("/courses/%s/students/%s", courseID, studentID)
	return doJSONRequest(t, app, "DELETE", path, nil, authHeaders(access))
}

func PostExamCreate(t *testing.T, app *fiber.App, access *HTTPAccess, body map[string]any) *HTTPResponse {
	return doJSONRequest(t, app, "POST", "/exams/", body, authHeaders(access))
}

func PostExamDefaultCodeTemplates(t *testing.T, app *fiber.App, access *HTTPAccess, body map[string]any) *HTTPResponse {
	return doJSONRequest(t, app, "POST", "/challenges/default-code-templates", body, authHeaders(access))
}

func PatchExamUpdate(t *testing.T, app *fiber.App, access *HTTPAccess, examID string, body map[string]any) *HTTPResponse {
	path := fmt.Sprintf("/exams/%s", examID)
	return doJSONRequest(t, app, "PATCH", path, body, authHeaders(access))
}

func GetExamByID(t *testing.T, app *fiber.App, access *HTTPAccess, examID string) *HTTPResponse {
	path := fmt.Sprintf("/exams/%s", examID)
	return doJSONRequest(t, app, "GET", path, nil, authHeaders(access))
}

func DeleteExamByID(t *testing.T, app *fiber.App, access *HTTPAccess, examID string) *HTTPResponse {
	path := fmt.Sprintf("/exams/%s", examID)
	return doJSONRequest(t, app, "DELETE", path, nil, authHeaders(access))
}

func PostExamClose(t *testing.T, app *fiber.App, access *HTTPAccess, examID string) *HTTPResponse {
	path := fmt.Sprintf("/exams/%s/close", examID)
	return doJSONRequest(t, app, "POST", path, nil, authHeaders(access))
}

func GetExamsByCourseID(t *testing.T, app *fiber.App, access *HTTPAccess, courseID string) *HTTPResponse {
	path := fmt.Sprintf("/exams/course/%s", courseID)
	return doJSONRequest(t, app, "GET", path, nil, authHeaders(access))
}

func GetPublicExams(t *testing.T, app *fiber.App, access *HTTPAccess) *HTTPResponse {
	return doJSONRequest(t, app, "GET", "/exams/public", nil, authHeaders(access))
}

func GetExamItems(t *testing.T, app *fiber.App, access *HTTPAccess, examID string) *HTTPResponse {
	path := fmt.Sprintf("/exams/%s/items", examID)
	return doJSONRequest(t, app, "GET", path, nil, authHeaders(access))
}

func PostChallengeCreate(t *testing.T, app *fiber.App, access *HTTPAccess, body map[string]any) *HTTPResponse {
	return doJSONRequest(t, app, "POST", "/challenges/", body, authHeaders(access))
}

func GetChallengeByID(t *testing.T, app *fiber.App, access *HTTPAccess, challengeID string) *HTTPResponse {
	path := fmt.Sprintf("/challenges/%s", challengeID)
	return doJSONRequest(t, app, "GET", path, nil, authHeaders(access))
}

func PatchChallengeUpdate(t *testing.T, app *fiber.App, access *HTTPAccess, challengeID string, body map[string]any) *HTTPResponse {
	path := fmt.Sprintf("/challenges/%s", challengeID)
	return doJSONRequest(t, app, "PATCH", path, body, authHeaders(access))
}

func PostChallengePublish(t *testing.T, app *fiber.App, access *HTTPAccess, challengeID string) *HTTPResponse {
	path := fmt.Sprintf("/challenges/%s/publish", challengeID)
	return doJSONRequest(t, app, "POST", path, nil, authHeaders(access))
}

func PostChallengeArchive(t *testing.T, app *fiber.App, access *HTTPAccess, challengeID string) *HTTPResponse {
	path := fmt.Sprintf("/challenges/%s/archive", challengeID)
	return doJSONRequest(t, app, "POST", path, nil, authHeaders(access))
}

func PostChallengeFork(t *testing.T, app *fiber.App, access *HTTPAccess, challengeID string) *HTTPResponse {
	path := fmt.Sprintf("/challenges/%s/fork", challengeID)
	return doJSONRequest(t, app, "POST", path, nil, authHeaders(access))
}

func GetPublicChallenges(t *testing.T, app *fiber.App, access *HTTPAccess) *HTTPResponse {
	return doJSONRequest(t, app, "GET", "/challenges/public", nil, authHeaders(access))
}

func DeleteChallengeByID(t *testing.T, app *fiber.App, access *HTTPAccess, challengeID string) *HTTPResponse {
	path := fmt.Sprintf("/challenges/%s", challengeID)
	return doJSONRequest(t, app, "DELETE", path, nil, authHeaders(access))
}

func PostTestCaseCreate(t *testing.T, app *fiber.App, access *HTTPAccess, body map[string]any) *HTTPResponse {
	return doJSONRequest(t, app, "POST", "/test-cases/", body, authHeaders(access))
}

func PatchTestCaseUpdate(t *testing.T, app *fiber.App, access *HTTPAccess, testCaseID string, body map[string]any) *HTTPResponse {
	path := fmt.Sprintf("/test-cases/%s", testCaseID)
	return doJSONRequest(t, app, "PATCH", path, body, authHeaders(access))
}

func GetTestCasesByChallenge(t *testing.T, app *fiber.App, access *HTTPAccess, challengeID, examID string) *HTTPResponse {
	path := fmt.Sprintf("/test-cases/challenge/%s", challengeID)
	if examID != "" {
		path = fmt.Sprintf("%s?exam_id=%s", path, examID)
	}
	return doJSONRequest(t, app, "GET", path, nil, authHeaders(access))
}

func DeleteTestCaseByID(t *testing.T, app *fiber.App, access *HTTPAccess, testCaseID string) *HTTPResponse {
	path := fmt.Sprintf("/test-cases/%s", testCaseID)
	return doJSONRequest(t, app, "DELETE", path, nil, authHeaders(access))
}

func PostExamItemCreate(t *testing.T, app *fiber.App, access *HTTPAccess, body map[string]any) *HTTPResponse {
	return doJSONRequest(t, app, "POST", "/exam-items/", body, authHeaders(access))
}

func PatchExamItemUpdate(t *testing.T, app *fiber.App, access *HTTPAccess, examItemID string, body map[string]any) *HTTPResponse {
	path := fmt.Sprintf("/exam-items/%s", examItemID)
	return doJSONRequest(t, app, "PATCH", path, body, authHeaders(access))
}

func DeleteExamItemByID(t *testing.T, app *fiber.App, access *HTTPAccess, examItemID string) *HTTPResponse {
	path := fmt.Sprintf("/exam-items/%s", examItemID)
	return doJSONRequest(t, app, "DELETE", path, nil, authHeaders(access))
}

func PostSessionCreate(t *testing.T, app *fiber.App, access *HTTPAccess, body map[string]any) *HTTPResponse {
	return doJSONRequest(t, app, "POST", "/submissions/sessions/", body, authHeaders(access))
}

func PostSessionHeartbeat(t *testing.T, app *fiber.App, access *HTTPAccess, sessionID string) *HTTPResponse {
	path := fmt.Sprintf("/submissions/sessions/%s/heartbeat", sessionID)
	return doJSONRequest(t, app, "POST", path, nil, authHeaders(access))
}

func PostSessionBlock(t *testing.T, app *fiber.App, access *HTTPAccess, sessionID string) *HTTPResponse {
	path := fmt.Sprintf("/submissions/sessions/%s/block", sessionID)
	return doJSONRequest(t, app, "POST", path, nil, authHeaders(access))
}

func PostSessionClose(t *testing.T, app *fiber.App, access *HTTPAccess, sessionID string) *HTTPResponse {
	path := fmt.Sprintf("/submissions/sessions/%s/close", sessionID)
	return doJSONRequest(t, app, "POST", path, nil, authHeaders(access))
}

func GetActiveSession(t *testing.T, app *fiber.App, access *HTTPAccess, userID string) *HTTPResponse {
	path := "/submissions/sessions/active"
	if userID != "" {
		path = fmt.Sprintf("/submissions/sessions/active?user_id=%s", userID)
	}
	return doJSONRequest(t, app, "GET", path, nil, authHeaders(access))
}

func PostSubmissionCreate(t *testing.T, app *fiber.App, access *HTTPAccess, body map[string]any) *HTTPResponse {
	return doJSONRequest(t, app, "POST", "/submissions/", body, authHeaders(access))
}

func PostSubmissionCreateWithoutScore(t *testing.T, app *fiber.App, access *HTTPAccess, body map[string]any) *HTTPResponse {
	return doJSONRequest(t, app, "POST", "/submissions/without-score", body, authHeaders(access))
}

func PostSubmissionCreateCustom(t *testing.T, app *fiber.App, access *HTTPAccess, body map[string]any) *HTTPResponse {
	return doJSONRequest(t, app, "POST", "/submissions/custom", body, authHeaders(access))
}

func GetSubmissionByID(t *testing.T, app *fiber.App, access *HTTPAccess, submissionID string) *HTTPResponse {
	path := fmt.Sprintf("/submissions/%s", submissionID)
	return doJSONRequest(t, app, "GET", path, nil, authHeaders(access))
}

func GetSubmissionsByChallenge(t *testing.T, app *fiber.App, access *HTTPAccess, challengeID string) *HTTPResponse {
	path := fmt.Sprintf("/submissions/challenge/%s", challengeID)
	return doJSONRequest(t, app, "GET", path, nil, authHeaders(access))
}

func GetSubmissionsByUser(t *testing.T, app *fiber.App, access *HTTPAccess, userID string) *HTTPResponse {
	path := fmt.Sprintf("/submissions/user/%s", userID)
	return doJSONRequest(t, app, "GET", path, nil, authHeaders(access))
}

func PatchSubmissionResult(t *testing.T, app *fiber.App, workerKey, resultID string, body map[string]any) *HTTPResponse {
	headers := map[string]string{"WorkerKey": workerKey}
	path := fmt.Sprintf("/submissions/results/%s", resultID)
	return doJSONRequest(t, app, "PATCH", path, body, headers)
}
