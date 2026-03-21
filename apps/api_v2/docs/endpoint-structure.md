# API v1 Endpoint Structure (Mapped to api_v2 mocks)

This document recognizes each original API v1 endpoint with:
- Address (method + path)
- Input (body/query/params)
- Output (mock JSON file used by api_v2)

## Auth
- `POST /auth/register`
  - Input: `{ username, password, role }`
  - Output: `test/http/mockup-responses/auth/post-register.json`
- `POST /auth/login`
  - Input: `{ username, password }`
  - Output: `test/http/mockup-responses/auth/post-login.json`
- `GET /auth/me`
  - Input: JWT auth header
  - Output: `test/http/mockup-responses/auth/get-me.json`

## AI
- `POST /ai/generate-challenge-ideas`
  - Input: `{ topic, difficulty?, count? }`
  - Output: `test/http/mockup-responses/ai/post-generate-challenge-ideas.json`
- `POST /ai/generate-test-cases`
  - Input: `{ challengeDescription, inputFormat, outputFormat, publicCount?, hiddenCount? }`
  - Output: `test/http/mockup-responses/ai/post-generate-test-cases.json`

## Challenges
- `POST /challenges`
  - Input: challenge payload with optional `publicTestCases` and `hiddenTestCases`
  - Output: `test/http/mockup-responses/challenges/post-create.json`
- `GET /challenges`
  - Input: none
  - Output: `test/http/mockup-responses/challenges/get-list.json`
- `GET /challenges/:id`
  - Input: path param `id`
  - Output: `test/http/mockup-responses/challenges/get-by-id.json`
- `PATCH /challenges/:id`
  - Input: path param `id` + challenge payload
  - Output: `test/http/mockup-responses/challenges/patch-update.json`
- `POST /challenges/:id/publish`
  - Input: path param `id`
  - Output: `test/http/mockup-responses/challenges/post-publish.json`
- `POST /challenges/:id/archive`
  - Input: path param `id`
  - Output: `test/http/mockup-responses/challenges/post-archive.json`

## Test Cases
- `POST /test-cases`
  - Input: `{ challengeId, name, input, expectedOutput, isSample?, points? }`
  - Output: `test/http/mockup-responses/test-cases/post-create.json`
- `GET /test-cases/challenge/:challengeId`
  - Input: path param `challengeId`, query `samplesOnly?`
  - Output: `test/http/mockup-responses/test-cases/get-by-challenge-id.json`
- `DELETE /test-cases/:id`
  - Input: path param `id`
  - Output: `test/http/mockup-responses/test-cases/delete-by-id.json`

## Courses
- `POST /courses/enroll`
  - Input: `{ enrollmentCode }`
  - Output: `test/http/mockup-responses/courses/post-enroll.json`
- `POST /courses`
  - Input: `{ name, code, period, groupNumber }`
  - Output: `test/http/mockup-responses/courses/post-create.json`
- `GET /courses/browse`
  - Input: none
  - Output: `test/http/mockup-responses/courses/get-browse.json`
- `GET /courses`
  - Input: none
  - Output: `test/http/mockup-responses/courses/get-list.json`
- `GET /courses/:id`
  - Input: path param `id`
  - Output: `test/http/mockup-responses/courses/get-by-id.json`
- `POST /courses/:id`
  - Input: path param `id` + `{ name, code, period, groupNumber }`
  - Output: `test/http/mockup-responses/courses/post-update.json`
- `POST /courses/:id/students`
  - Input: path param `id` + `{ studentId }`
  - Output: `test/http/mockup-responses/courses/post-add-student.json`
- `DELETE /courses/:id/students/:studentId`
  - Input: path params `id`, `studentId`
  - Output: `test/http/mockup-responses/courses/delete-student.json`
- `POST /courses/:id/challenges`
  - Input: path param `id` + `{ challengeId }`
  - Output: `test/http/mockup-responses/courses/post-assign-challenge.json`
- `GET /courses/:id/students`
  - Input: path param `id`
  - Output: `test/http/mockup-responses/courses/get-students.json`
- `GET /courses/:id/challenges`
  - Input: path param `id`
  - Output: `test/http/mockup-responses/courses/get-challenges.json`

## Exams
- `POST /exams`
  - Input: `{ title, description?, courseId, startTime, endTime, durationMinutes, challenges[] }`
  - Output: `test/http/mockup-responses/exams/post-create.json`
- `GET /exams/course/:courseId`
  - Input: path param `courseId`
  - Output: `test/http/mockup-responses/exams/get-by-course-id.json`
- `GET /exams/:id`
  - Input: path param `id`
  - Output: `test/http/mockup-responses/exams/get-by-id.json`

## Submissions
- `POST /submissions`
  - Input: `{ challengeId, code, language, examId? }`
  - Output: `test/http/mockup-responses/submissions/post-create.json`
- `GET /submissions/:id`
  - Input: path param `id`
  - Output: `test/http/mockup-responses/submissions/get-by-id.json`
- `GET /submissions`
  - Input: query `challengeId?`, `status?`, `limit?`, `offset?`
  - Output: `test/http/mockup-responses/submissions/get-list.json`

## Leaderboard
- `GET /leaderboard/challenge/:id`
  - Input: path param `id`
  - Output: `test/http/mockup-responses/leaderboard/get-challenge-id.json`
- `GET /leaderboard/course/:id`
  - Input: path param `id`
  - Output: `test/http/mockup-responses/leaderboard/get-course-id.json`

## Observability and Health
- `GET /metrics`
  - Input: none
  - Output: `test/http/mockup-responses/metrics/get.json`
- `GET /health`
  - Input: none
  - Output: `test/http/mockup-responses/health/get.json`
- `GET /cache/health`
  - Input: none
  - Output: `test/http/mockup-responses/cache/get-health.json`
- `GET /db/health`
  - Input: none
  - Output: `test/http/mockup-responses/db/get-health.json`
