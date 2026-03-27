package utils

import (
	"fmt"
)

// --- AUTH ---
func PostAuthRegister(headers map[string]string, body any) (int, []byte, error) {
	return callEndpoint("POST", "/auth/register", nil, nil, headers, body)
}
func PostAuthLogin(headers map[string]string, body any) (int, []byte, error) {
	return callEndpoint("POST", "/auth/login", nil, nil, headers, body)
}
func PostAuthRefreshToken(headers map[string]string, body any) (int, []byte, error) {
	return callEndpoint("POST", "/auth/refresh-token", nil, nil, headers, body)
}
func GetAuthMe(headers map[string]string) (int, []byte, error) {
	return callEndpoint("GET", "/auth/me", nil, nil, headers, nil)
}

// --- AI ---
func PostAiGenerateChallengeIdeas(headers map[string]string, body any) (int, []byte, error) {
	return callEndpoint("POST", "/ai/generate-challenge-ideas", nil, nil, headers, body)
}
func PostAiGenerateTestCases(headers map[string]string, body any) (int, []byte, error) {
	return callEndpoint("POST", "/ai/generate-test-cases", nil, nil, headers, body)
}

// --- CHALLENGES ---
func PostChallenges(headers map[string]string, body any) (int, []byte, error) {
	return callEndpoint("POST", "/challenges/", nil, nil, headers, body)
}
func GetChallenges(headers map[string]string, query map[string]any) (int, []byte, error) {
	return callEndpoint("GET", "/challenges/", nil, query, headers, nil)
}
func GetChallengesById(headers map[string]string, path map[string]any) (int, []byte, error) {
	return callEndpoint("GET", "/challenges/{id}", path, nil, headers, nil)
}
func PatchChallengesById(headers map[string]string, path map[string]any, body any) (int, []byte, error) {
	return callEndpoint("PATCH", "/challenges/{id}", path, nil, headers, body)
}
func DeleteChallengesById(headers map[string]string, path map[string]any) (int, []byte, error) {
	return callEndpoint("DELETE", "/challenges/{id}", path, nil, headers, nil)
}
func PostChallengesPublish(headers map[string]string, path map[string]any) (int, []byte, error) {
	return callEndpoint("POST", "/challenges/{id}/publish", path, nil, headers, nil)
}
func PostChallengesArchive(headers map[string]string, path map[string]any) (int, []byte, error) {
	return callEndpoint("POST", "/challenges/{id}/archive", path, nil, headers, nil)
}

// --- TEST CASES ---
func PostTestCases(headers map[string]string, body any) (int, []byte, error) {
	return callEndpoint("POST", "/test-cases/", nil, nil, headers, body)
}
func GetTestCasesByChallenge(headers map[string]string, path map[string]any, query map[string]any) (int, []byte, error) {
	return callEndpoint("GET", "/test-cases/challenge/{id}", path, query, headers, nil)
}
func PatchTestCasesById(headers map[string]string, path map[string]any, body any) (int, []byte, error) {
	return callEndpoint("PATCH", "/test-cases/{id}", path, nil, headers, body)
}
func DeleteTestCasesById(headers map[string]string, path map[string]any) (int, []byte, error) {
	return callEndpoint("DELETE", "/test-cases/{id}", path, nil, headers, nil)
}

// --- COURSES ---
func PostCoursesEnroll(headers map[string]string, body any) (int, []byte, error) {
	return callEndpoint("POST", "/courses/enroll", nil, nil, headers, body)
}
func PostCourses(headers map[string]string, body any) (int, []byte, error) {
	return callEndpoint("POST", "/courses/", nil, nil, headers, body)
}
func GetCourses(headers map[string]string, query map[string]any) (int, []byte, error) {
	return callEndpoint("GET", "/courses/", nil, query, headers, nil)
}
func GetCoursesById(headers map[string]string, path map[string]any) (int, []byte, error) {
	return callEndpoint("GET", "/courses/{id}", path, nil, headers, nil)
}
func PostCoursesById(headers map[string]string, path map[string]any, body any) (int, []byte, error) {
	return callEndpoint("POST", "/courses/{id}", path, nil, headers, body)
}
func DeleteCoursesById(headers map[string]string, path map[string]any) (int, []byte, error) {
	return callEndpoint("DELETE", "/courses/{id}", path, nil, headers, nil)
}
func PostCoursesAddStudent(headers map[string]string, path map[string]any, body any) (int, []byte, error) {
	return callEndpoint("POST", "/courses/{id}/students", path, nil, headers, body)
}
func DeleteCoursesStudent(headers map[string]string, path map[string]any) (int, []byte, error) {
	return callEndpoint("DELETE", "/courses/{id}/students/{student_id}", path, nil, headers, nil)
}
func GetCoursesStudents(headers map[string]string, path map[string]any) (int, []byte, error) {
	return callEndpoint("GET", "/courses/{id}/students", path, nil, headers, nil)
}

// --- EXAMS ---
func PostExams(headers map[string]string, body any) (int, []byte, error) {
	return callEndpoint("POST", "/exams/", nil, nil, headers, body)
}
func GetExamsByCourse(headers map[string]string, path map[string]any) (int, []byte, error) {
	return callEndpoint("GET", "/exams/course/{course_id}", path, nil, headers, nil)
}
func GetExamsById(headers map[string]string, path map[string]any) (int, []byte, error) {
	return callEndpoint("GET", "/exams/{id}", path, nil, headers, nil)
}
func PatchExamsById(headers map[string]string, path map[string]any, body any) (int, []byte, error) {
	return callEndpoint("PATCH", "/exams/{id}", path, nil, headers, body)
}
func DeleteExamsById(headers map[string]string, path map[string]any) (int, []byte, error) {
	return callEndpoint("DELETE", "/exams/{id}", path, nil, headers, nil)
}
func PostExamsVisibility(headers map[string]string, path map[string]any, body any) (int, []byte, error) {
	return callEndpoint("POST", "/exams/{id}/visibility", path, nil, headers, body)
}
func PostExamsClose(headers map[string]string, path map[string]any) (int, []byte, error) {
	return callEndpoint("POST", "/exams/{id}/close", path, nil, headers, nil)
}
func GetExamItems(headers map[string]string, path map[string]any) (int, []byte, error) {
	return callEndpoint("GET", "/exams/{exam_id}/items", path, nil, headers, nil)
}

// --- EXAM ITEMS ---
func PostExamItems(headers map[string]string, body any) (int, []byte, error) {
	return callEndpoint("POST", "/exam-items/", nil, nil, headers, body)
}
func PatchExamItemsById(headers map[string]string, path map[string]any, body any) (int, []byte, error) {
	return callEndpoint("PATCH", "/exam-items/{id}", path, nil, headers, body)
}
func DeleteExamItemsById(headers map[string]string, path map[string]any) (int, []byte, error) {
	return callEndpoint("DELETE", "/exam-items/{id}", path, nil, headers, nil)
}

// --- SUBMISSIONS ---
func PostSubmissions(headers map[string]string, body any) (int, []byte, error) {
	return callEndpoint("POST", "/submissions/", nil, nil, headers, body)
}
func PatchSubmissionsResult(headers map[string]string, path map[string]any, body any) (int, []byte, error) {
	return callEndpoint("PATCH", "/submissions/results/{result_id}", path, nil, headers, body)
}
func GetSubmissionsUser(headers map[string]string, path map[string]any, query map[string]any) (int, []byte, error) {
	return callEndpoint("GET", "/submissions/user/{user_id}", path, query, headers, nil)
}
func GetSubmissionsSession(headers map[string]string, path map[string]any, query map[string]any) (int, []byte, error) {
	return callEndpoint("GET", "/submissions/session/{session_id}", path, query, headers, nil)
}
func GetSubmissionsChallenge(headers map[string]string, path map[string]any, query map[string]any) (int, []byte, error) {
	return callEndpoint("GET", "/submissions/challenge/{challenge_id}", path, query, headers, nil)
}
func GetSubmissionsById(headers map[string]string, path map[string]any) (int, []byte, error) {
	return callEndpoint("GET", "/submissions/{id}", path, nil, headers, nil)
}
func PostSubmissionsSessions(headers map[string]string, body any) (int, []byte, error) {
	return callEndpoint("POST", "/submissions/sessions/", nil, nil, headers, body)
}
func GetSubmissionsSessionsById(headers map[string]string, path map[string]any) (int, []byte, error) {
	return callEndpoint("GET", "/submissions/sessions/{id}", path, nil, headers, nil)
}
func PostSubmissionsSessionsHeartbeat(headers map[string]string, path map[string]any) (int, []byte, error) {
	return callEndpoint("POST", "/submissions/sessions/{id}/heartbeat", path, nil, headers, nil)
}

// --- LEADERBOARD ---
func GetLeaderboardChallenge(headers map[string]string, path map[string]any) (int, []byte, error) {
	return callEndpoint("GET", "/leaderboard/challenge/{id}", path, nil, headers, nil)
}
func GetLeaderboardCourse(headers map[string]string, path map[string]any) (int, []byte, error) {
	return callEndpoint("GET", "/leaderboard/course/{id}", path, nil, headers, nil)
}

// --- METRICS, HEALTH, CACHE, DB ---
func GetMetrics(headers map[string]string) (int, []byte, error) {
	return callEndpoint("GET", "/metrics", nil, nil, headers, nil)
}
func GetHealth(headers map[string]string) (int, []byte, error) {
	return callEndpoint("GET", "/health", nil, nil, headers, nil)
}
func GetCacheHealth(headers map[string]string) (int, []byte, error) {
	return callEndpoint("GET", "/cache/health", nil, nil, headers, nil)
}
func GetDbHealth(headers map[string]string) (int, []byte, error) {
	return callEndpoint("GET", "/db/health", nil, nil, headers, nil)
}

// --- CORE REQUEST FUNCTION ---
func callEndpoint(method, path string, pathParams, queryParams map[string]any, headers map[string]string, body any) (int, []byte, error) {
	app, err := InitApp()
	if err != nil {
		return -1, nil, fmt.Errorf("init app: %w", err)
	}
	fullPath := buildRequestPath(path, pathParams, queryParams)
	return DoJSONRequest(app, method, fullPath, body, headers)
}
