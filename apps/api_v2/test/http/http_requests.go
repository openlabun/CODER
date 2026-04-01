package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type HTTPRequestDTO struct {
	Status       int           `json:"status"`
	BodyResponse []byte        `json:"body_response"`
	Time         time.Duration `json:"time"`
}

type AuthInput struct {
	Token     string
	UserEmail string
}

func (a AuthInput) headers() map[string]string {
	headers := map[string]string{}
	if a.Token != "" {
		headers["Authorization"] = "Bearer " + a.Token
	}
	if a.UserEmail != "" {
		headers["X-User-Email"] = a.UserEmail
	}

	return headers
}

func performEndpointRequest(method, path string, pathParams, queryParams map[string]any, auth AuthInput, body any) (HTTPRequestDTO, error) {
	app, err := InitApp()
	if err != nil {
		return HTTPRequestDTO{}, fmt.Errorf("init app: %w", err)
	}

	fullPath := buildRequestPath(path, pathParams, queryParams)
	start := time.Now()
	status, responseBody, err := DoJSONRequest(app, method, fullPath, body, auth.headers())
	elapsed := time.Since(start)
	dto := HTTPRequestDTO{Status: status, BodyResponse: responseBody, Time: elapsed}
	if err != nil {
		return dto, err
	}

	if status < 200 || status >= 300 {
		return dto, fmt.Errorf("unexpected status=%d path=%s body=%s", status, fullPath, string(responseBody))
	}

	trimmed := bytes.TrimSpace(responseBody)
	if len(trimmed) > 0 && (trimmed[0] == '{' || trimmed[0] == '[') {
		var payload any
		if unmarshalErr := json.Unmarshal(trimmed, &payload); unmarshalErr != nil {
			return dto, fmt.Errorf("invalid json response for path=%s: %w", fullPath, unmarshalErr)
		}
	}

	return dto, nil
}

// --- DOCS ---
func ReqGetDocsOpenAPI() (HTTPRequestDTO, error) {
	return performEndpointRequest("GET", "/docs/openapi.yaml", nil, nil, AuthInput{}, nil)
}

func ReqGetDocs() (HTTPRequestDTO, error) {
	return performEndpointRequest("GET", "/docs", nil, nil, AuthInput{}, nil)
}

// --- AUTH ---
func ReqPostAuthRegister(email, name, password string) (HTTPRequestDTO, error) {
	body := map[string]any{"email": email, "name": name, "password": password}
	return performEndpointRequest("POST", "/auth/register", nil, nil, AuthInput{}, body)
}

func ReqPostAuthLogin(email, password string) (HTTPRequestDTO, error) {
	body := map[string]any{"email": email, "password": password}
	return performEndpointRequest("POST", "/auth/login", nil, nil, AuthInput{}, body)
}

func ReqPostAuthRefreshToken(refreshToken string) (HTTPRequestDTO, error) {
	body := map[string]any{"refresh_token": refreshToken}
	return performEndpointRequest("POST", "/auth/refresh-token", nil, nil, AuthInput{}, body)
}

func ReqGetAuthMe(token, userEmail string) (HTTPRequestDTO, error) {
	return performEndpointRequest("GET", "/auth/me", nil, nil, AuthInput{Token: token, UserEmail: userEmail}, nil)
}

// --- AI ---
func ReqPostAiGenerateChallengeIdeas(token, userEmail string, body any) (HTTPRequestDTO, error) {
	return performEndpointRequest("POST", "/ai/generate-challenge-ideas", nil, nil, AuthInput{Token: token, UserEmail: userEmail}, body)
}

func ReqPostAiGenerateTestCases(token, userEmail string, body any) (HTTPRequestDTO, error) {
	return performEndpointRequest("POST", "/ai/generate-test-cases", nil, nil, AuthInput{Token: token, UserEmail: userEmail}, body)
}

func ReqPostAiGenerateFullChallenge(token, userEmail string, body any) (HTTPRequestDTO, error) {
	return performEndpointRequest("POST", "/ai/generate-full-challenge", nil, nil, AuthInput{Token: token, UserEmail: userEmail}, body)
}

func ReqPostAiGenerateExam(token, userEmail string, body any) (HTTPRequestDTO, error) {
	return performEndpointRequest("POST", "/ai/generate-exam", nil, nil, AuthInput{Token: token, UserEmail: userEmail}, body)
}

// --- CHALLENGES ---
func ReqPostChallenges(token, userEmail string, body any) (HTTPRequestDTO, error) {
	return performEndpointRequest("POST", "/challenges/", nil, nil, AuthInput{Token: token, UserEmail: userEmail}, body)
}

func ReqGetChallenges(token, userEmail string, examID *string) (HTTPRequestDTO, error) {
	query := map[string]any{}
	if examID != nil {
		query["examId"] = *examID
	}
	return performEndpointRequest("GET", "/challenges/", nil, query, AuthInput{Token: token, UserEmail: userEmail}, nil)
}

func ReqGetChallengesPublic(token, userEmail string, tag, difficulty *string) (HTTPRequestDTO, error) {
	query := map[string]any{}
	if tag != nil {
		query["tag"] = *tag
	}
	if difficulty != nil {
		query["difficulty"] = *difficulty
	}
	return performEndpointRequest("GET", "/challenges/public", nil, query, AuthInput{Token: token, UserEmail: userEmail}, nil)
}

func ReqGetChallengeByID(token, userEmail, id string) (HTTPRequestDTO, error) {
	return performEndpointRequest("GET", "/challenges/{id}", map[string]any{"id": id}, nil, AuthInput{Token: token, UserEmail: userEmail}, nil)
}

func ReqPatchChallengeByID(token, userEmail, id string, body any) (HTTPRequestDTO, error) {
	return performEndpointRequest("PATCH", "/challenges/{id}", map[string]any{"id": id}, nil, AuthInput{Token: token, UserEmail: userEmail}, body)
}

func ReqDeleteChallengeByID(token, userEmail, id string) (HTTPRequestDTO, error) {
	return performEndpointRequest("DELETE", "/challenges/{id}", map[string]any{"id": id}, nil, AuthInput{Token: token, UserEmail: userEmail}, nil)
}

func ReqPostChallengePublish(token, userEmail, id string) (HTTPRequestDTO, error) {
	return performEndpointRequest("POST", "/challenges/{id}/publish", map[string]any{"id": id}, nil, AuthInput{Token: token, UserEmail: userEmail}, nil)
}

func ReqPostChallengeArchive(token, userEmail, id string) (HTTPRequestDTO, error) {
	return performEndpointRequest("POST", "/challenges/{id}/archive", map[string]any{"id": id}, nil, AuthInput{Token: token, UserEmail: userEmail}, nil)
}

func ReqPostChallengeFork(token, userEmail, id string) (HTTPRequestDTO, error) {
	return performEndpointRequest("POST", "/challenges/{id}/fork", map[string]any{"id": id}, nil, AuthInput{Token: token, UserEmail: userEmail}, nil)
}

// --- TEST CASES ---
func ReqPostTestCases(token, userEmail string, body any) (HTTPRequestDTO, error) {
	return performEndpointRequest("POST", "/test-cases/", nil, nil, AuthInput{Token: token, UserEmail: userEmail}, body)
}

func ReqGetTestCasesByChallengeID(token, userEmail, challengeID string, examID *string) (HTTPRequestDTO, error) {
	query := map[string]any{}
	if examID != nil {
		query["exam_id"] = *examID
	}
	return performEndpointRequest("GET", "/test-cases/challenge/{challengeId}", map[string]any{"challengeId": challengeID}, query, AuthInput{Token: token, UserEmail: userEmail}, nil)
}

func ReqPatchTestCaseByID(token, userEmail, id string, body any) (HTTPRequestDTO, error) {
	return performEndpointRequest("PATCH", "/test-cases/{id}", map[string]any{"id": id}, nil, AuthInput{Token: token, UserEmail: userEmail}, body)
}

func ReqDeleteTestCaseByID(token, userEmail, id string) (HTTPRequestDTO, error) {
	return performEndpointRequest("DELETE", "/test-cases/{id}", map[string]any{"id": id}, nil, AuthInput{Token: token, UserEmail: userEmail}, nil)
}

// --- COURSES ---
func ReqPostCoursesEnroll(token, userEmail, courseID, studentID string) (HTTPRequestDTO, error) {
	body := map[string]any{"course_id": courseID, "student_id": studentID}
	return performEndpointRequest("POST", "/courses/enroll", nil, nil, AuthInput{Token: token, UserEmail: userEmail}, body)
}

func ReqPostCourses(token, userEmail string, body any) (HTTPRequestDTO, error) {
	return performEndpointRequest("POST", "/courses/", nil, nil, AuthInput{Token: token, UserEmail: userEmail}, body)
}

func ReqGetCourses(token, userEmail string, query map[string]any) (HTTPRequestDTO, error) {
	return performEndpointRequest("GET", "/courses/", nil, query, AuthInput{Token: token, UserEmail: userEmail}, nil)
}

func ReqGetCourseByID(token, userEmail, id string) (HTTPRequestDTO, error) {
	return performEndpointRequest("GET", "/courses/{id}", map[string]any{"id": id}, nil, AuthInput{Token: token, UserEmail: userEmail}, nil)
}

func ReqPostCourseByID(token, userEmail, id string, body any) (HTTPRequestDTO, error) {
	return performEndpointRequest("POST", "/courses/{id}", map[string]any{"id": id}, nil, AuthInput{Token: token, UserEmail: userEmail}, body)
}

func ReqDeleteCourseByID(token, userEmail, id string) (HTTPRequestDTO, error) {
	return performEndpointRequest("DELETE", "/courses/{id}", map[string]any{"id": id}, nil, AuthInput{Token: token, UserEmail: userEmail}, nil)
}

func ReqPostCourseAddStudent(token, userEmail, id string, body any) (HTTPRequestDTO, error) {
	return performEndpointRequest("POST", "/courses/{id}/students", map[string]any{"id": id}, nil, AuthInput{Token: token, UserEmail: userEmail}, body)
}

func ReqGetCourseStudents(token, userEmail, id string) (HTTPRequestDTO, error) {
	return performEndpointRequest("GET", "/courses/{id}/students", map[string]any{"id": id}, nil, AuthInput{Token: token, UserEmail: userEmail}, nil)
}

func ReqDeleteCourseStudent(token, userEmail, id, studentID string) (HTTPRequestDTO, error) {
	return performEndpointRequest("DELETE", "/courses/{id}/students/{studentId}", map[string]any{"id": id, "studentId": studentID}, nil, AuthInput{Token: token, UserEmail: userEmail}, nil)
}

// --- EXAMS ---
func ReqPostExams(token, userEmail string, body any) (HTTPRequestDTO, error) {
	return performEndpointRequest("POST", "/exams/", nil, nil, AuthInput{Token: token, UserEmail: userEmail}, body)
}

func ReqGetExamsByCourseID(token, userEmail, courseID string) (HTTPRequestDTO, error) {
	return performEndpointRequest("GET", "/exams/course/{courseId}", map[string]any{"courseId": courseID}, nil, AuthInput{Token: token, UserEmail: userEmail}, nil)
}

func ReqGetExamsPublic(token, userEmail string) (HTTPRequestDTO, error) {
	return performEndpointRequest("GET", "/exams/public", nil, nil, AuthInput{Token: token, UserEmail: userEmail}, nil)
}

func ReqGetExamByID(token, userEmail, id string) (HTTPRequestDTO, error) {
	return performEndpointRequest("GET", "/exams/{id}", map[string]any{"id": id}, nil, AuthInput{Token: token, UserEmail: userEmail}, nil)
}

func ReqPatchExamByID(token, userEmail, id string, body any) (HTTPRequestDTO, error) {
	return performEndpointRequest("PATCH", "/exams/{id}", map[string]any{"id": id}, nil, AuthInput{Token: token, UserEmail: userEmail}, body)
}

func ReqDeleteExamByID(token, userEmail, id string) (HTTPRequestDTO, error) {
	return performEndpointRequest("DELETE", "/exams/{id}", map[string]any{"id": id}, nil, AuthInput{Token: token, UserEmail: userEmail}, nil)
}

func ReqPostExamVisibility(token, userEmail, id string, body any) (HTTPRequestDTO, error) {
	return performEndpointRequest("POST", "/exams/{id}/visibility", map[string]any{"id": id}, nil, AuthInput{Token: token, UserEmail: userEmail}, body)
}

func ReqPostExamClose(token, userEmail, id string) (HTTPRequestDTO, error) {
	return performEndpointRequest("POST", "/exams/{id}/close", map[string]any{"id": id}, nil, AuthInput{Token: token, UserEmail: userEmail}, nil)
}

func ReqGetExamItems(token, userEmail, examID string) (HTTPRequestDTO, error) {
	return performEndpointRequest("GET", "/exams/{examId}/items", map[string]any{"examId": examID}, nil, AuthInput{Token: token, UserEmail: userEmail}, nil)
}

// --- EXAM ITEMS ---
func ReqPostExamItems(token, userEmail string, body any) (HTTPRequestDTO, error) {
	return performEndpointRequest("POST", "/exam-items/", nil, nil, AuthInput{Token: token, UserEmail: userEmail}, body)
}

func ReqPatchExamItemByID(token, userEmail, id string, body any) (HTTPRequestDTO, error) {
	return performEndpointRequest("PATCH", "/exam-items/{id}", map[string]any{"id": id}, nil, AuthInput{Token: token, UserEmail: userEmail}, body)
}

func ReqDeleteExamItemByID(token, userEmail, id string) (HTTPRequestDTO, error) {
	return performEndpointRequest("DELETE", "/exam-items/{id}", map[string]any{"id": id}, nil, AuthInput{Token: token, UserEmail: userEmail}, nil)
}

// --- SUBMISSIONS ---
func ReqPostSubmissions(token, userEmail string, body any) (HTTPRequestDTO, error) {
	return performEndpointRequest("POST", "/submissions/", nil, nil, AuthInput{Token: token, UserEmail: userEmail}, body)
}

func ReqPatchSubmissionResult(token, userEmail, resultID string, body any) (HTTPRequestDTO, error) {
	return performEndpointRequest("PATCH", "/submissions/results/{resultId}", map[string]any{"resultId": resultID}, nil, AuthInput{Token: token, UserEmail: userEmail}, body)
}

func ReqGetSubmissionsByUser(token, userEmail, userID string, query map[string]any) (HTTPRequestDTO, error) {
	return performEndpointRequest("GET", "/submissions/user/{userId}", map[string]any{"userId": userID}, query, AuthInput{Token: token, UserEmail: userEmail}, nil)
}

func ReqGetSubmissionsByChallenge(token, userEmail, challengeID string, query map[string]any) (HTTPRequestDTO, error) {
	return performEndpointRequest("GET", "/submissions/challenge/{challengeId}", map[string]any{"challengeId": challengeID}, query, AuthInput{Token: token, UserEmail: userEmail}, nil)
}

func ReqGetSubmissionByID(token, userEmail, id string) (HTTPRequestDTO, error) {
	return performEndpointRequest("GET", "/submissions/{id}", map[string]any{"id": id}, nil, AuthInput{Token: token, UserEmail: userEmail}, nil)
}

func ReqPostSubmissionSession(token, userEmail, userID, examID string) (HTTPRequestDTO, error) {
	body := map[string]any{"user_id": userID, "exam_id": examID}
	return performEndpointRequest("POST", "/submissions/sessions/", nil, nil, AuthInput{Token: token, UserEmail: userEmail}, body)
}

func ReqGetActiveSubmissionSession(token, userEmail string, targetUserID *string) (HTTPRequestDTO, error) {
	query := map[string]any{}
	if targetUserID != nil {
		query["user_id"] = *targetUserID
	}
	return performEndpointRequest("GET", "/submissions/sessions/active", nil, query, AuthInput{Token: token, UserEmail: userEmail}, nil)
}

func ReqPostSubmissionSessionHeartbeat(token, userEmail, id string) (HTTPRequestDTO, error) {
	return performEndpointRequest("POST", "/submissions/sessions/{id}/heartbeat", map[string]any{"id": id}, nil, AuthInput{Token: token, UserEmail: userEmail}, nil)
}

func ReqPostSubmissionSessionBlock(token, userEmail, id string) (HTTPRequestDTO, error) {
	return performEndpointRequest("POST", "/submissions/sessions/{id}/block", map[string]any{"id": id}, nil, AuthInput{Token: token, UserEmail: userEmail}, nil)
}

func ReqPostSubmissionSessionClose(token, userEmail, id string) (HTTPRequestDTO, error) {
	return performEndpointRequest("POST", "/submissions/sessions/{id}/close", map[string]any{"id": id}, nil, AuthInput{Token: token, UserEmail: userEmail}, nil)
}

// --- LEADERBOARD ---
func ReqGetLeaderboardChallenge(token, userEmail, id string) (HTTPRequestDTO, error) {
	return performEndpointRequest("GET", "/leaderboard/challenge/{id}", map[string]any{"id": id}, nil, AuthInput{Token: token, UserEmail: userEmail}, nil)
}

func ReqGetLeaderboardCourse(token, userEmail, id string) (HTTPRequestDTO, error) {
	return performEndpointRequest("GET", "/leaderboard/course/{id}", map[string]any{"id": id}, nil, AuthInput{Token: token, UserEmail: userEmail}, nil)
}

// --- METRICS, HEALTH, CACHE, DB ---
func ReqGetMetrics(token, userEmail string) (HTTPRequestDTO, error) {
	return performEndpointRequest("GET", "/metrics", nil, nil, AuthInput{Token: token, UserEmail: userEmail}, nil)
}

func ReqGetHealth(token, userEmail string) (HTTPRequestDTO, error) {
	return performEndpointRequest("GET", "/health", nil, nil, AuthInput{Token: token, UserEmail: userEmail}, nil)
}

func ReqGetCacheHealth(token, userEmail string) (HTTPRequestDTO, error) {
	return performEndpointRequest("GET", "/cache/health", nil, nil, AuthInput{Token: token, UserEmail: userEmail}, nil)
}

func ReqGetDBHealth(token, userEmail string) (HTTPRequestDTO, error) {
	return performEndpointRequest("GET", "/db/health", nil, nil, AuthInput{Token: token, UserEmail: userEmail}, nil)
}
