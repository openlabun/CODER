package submission_test

import (
	"fmt"
	"testing"
	"time"

	container "github.com/openlabun/CODER/apps/api_v2/internal/application/container"
	course_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/course"
	exam_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	submission_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"
	submission_usecases "github.com/openlabun/CODER/apps/api_v2/internal/application/usecases/submission"

	course_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"
	exam_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	submission_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"

	utils "github.com/openlabun/CODER/apps/api_v2/test/use_cases"
)

func TestCreateSession(t *testing.T) {
	t.Log("[STEP 1] Initialize application container with dependencies")
	app, err := container.BuildApplicationContainer()
	if err != nil {
		t.Fatalf("failed to build application: %v", err)
	}
	t.Log("[OK] Application initialized")

	t.Log("[STEP 2] Login professor user and create context with credentials")
	teacherAccess := utils.EnsureTeacherAccess(t, app)
	teacherCtx := utils.TeacherCourseCtx(teacherAccess)
	teacherID := teacherAccess.UserData.ID
	t.Logf("[OK] Teacher login successful. teacherID=%s", teacherID)

	now := time.Now().UTC()
	enrollmentCode := fmt.Sprintf("ENR-UC-SES-%d", now.UnixNano())
	courseCode := fmt.Sprintf("UC-SES-%d", now.Unix()%100000)

	var courseID string
	var examID string
	var studentSessionID string
	var teacherSessionID string

	defer func() {
		if studentSessionID != "" {
			t.Logf("[CLEANUP] Deleting student session %s", studentSessionID)
			_ = app.Dependencies.SessionRepository.DeleteSession(teacherCtx, studentSessionID)
		}
		if teacherSessionID != "" {
			t.Logf("[CLEANUP] Deleting teacher session %s", teacherSessionID)
			_ = app.Dependencies.SessionRepository.DeleteSession(teacherCtx, teacherSessionID)
		}
		if examID != "" {
			t.Logf("[CLEANUP] Deleting exam %s", examID)
			_, _ = app.ExamModule.DeleteExam.Execute(teacherCtx, exam_dtos.DeleteExamInput{ExamID: examID})
		}
		if courseID != "" {
			t.Logf("[CLEANUP] Deleting course %s", courseID)
			_ = app.CourseModule.DeleteCourse.Execute(teacherCtx, course_dtos.DeleteCourseInput{CourseID: courseID})
		}
	}()

	t.Log("[STEP 3] Create a course for testing")
	createdCourse, err := app.CourseModule.CreateCourse.Execute(teacherCtx, course_dtos.CreateCourseInput{
		Name:           "UC Submission Session Course",
		Description:    "Course for session use case test",
		Visibility:     string(course_entities.CourseVisibilityPublic),
		VisualIdentity: string(course_entities.CourseColourBlue),
		Code:           courseCode,
		Year:           now.Year(),
		Semester:       string(course_entities.AcademicFirstPeriod),
		EnrollmentCode: enrollmentCode,
		EnrollmentURL:  "https://example.test/enroll/" + enrollmentCode,
		TeacherID:      teacherID,
	})
	if err != nil {
		t.Fatalf("create course failed: %v", err)
	}
	if createdCourse == nil || createdCourse.ID == "" {
		t.Fatal("expected created course with ID")
	}
	courseID = createdCourse.ID
	t.Logf("[OK] Course created. courseID=%s", courseID)

	t.Log("[STEP 4] Create an exam for testing")
	startTime := now.Add(2 * time.Hour)
	endTime := startTime.Add(90 * time.Minute)
	endTimeStr := endTime.Format(time.RFC3339)

	createdExam, err := app.ExamModule.CreateExam.Execute(teacherCtx, exam_dtos.CreateExamInput{
		CourseID:             courseID,
		Title:                "UC Submission Session Exam",
		Description:          "Exam for session use case test",
		Visibility:           string(exam_entities.VisibilityCourse),
		StartTime:            startTime.Format(time.RFC3339),
		EndTime:              &endTimeStr,
		AllowLateSubmissions: false,
		TimeLimit:            3600,
		TryLimit:             2,
		ProfessorID:          teacherID,
	})
	if err != nil {
		t.Fatalf("create exam failed: %v", err)
	}
	if createdExam == nil || createdExam.ID == "" {
		t.Fatal("expected created exam with ID")
	}
	examID = createdExam.ID
	t.Logf("[OK] Exam created. examID=%s", examID)

	t.Log("[STEP 5] Login student user and create context with credentials")
	studentAccess := utils.EnsureStudentAccess(t, app)
	studentCtx := utils.StudentCtx(studentAccess)
	studentID := studentAccess.UserData.ID
	t.Logf("[OK] Student access ready. studentID=%s", studentID)

	t.Log("[STEP 6] Enroll student in course")
	_, err = app.CourseModule.EnrollInCourse.Execute(studentCtx, course_dtos.EnrolledInCourseInput{
		CourseID:  courseID,
		StudentID: &studentID,
		StudentEmail: nil,
	})
	if err != nil {
		t.Fatalf("enroll student failed: %v", err)
	}
	t.Log("[OK] Student enrolled in course")

	t.Log("[STEP 7] Create Session for Student and Exam")
	studentSession, err := app.SessionModule.CreateSession.Execute(studentCtx, submission_dtos.CreateSessionInput{
		UserID: studentID,
		ExamID: examID,
	})
	if err != nil {
		t.Fatalf("create session for student failed: %v", err)
	}
	if studentSession == nil || studentSession.ID == "" {
		t.Fatal("expected created session with ID")
	}
	studentSessionID = studentSession.ID
	t.Logf("[OK] Session created. sessionID=%s", studentSessionID)

	t.Log("[STEP 8] Assert session is created and returned successfully")
	if studentSession.StudentID != studentID {
		t.Fatalf("expected session studentID=%s, got=%s", studentID, studentSession.StudentID)
	}
	if studentSession.ExamID != examID {
		t.Fatalf("expected session examID=%s, got=%s", examID, studentSession.ExamID)
	}
	t.Log("[OK] Session payload validated")

	t.Log("[STEP 9] Try to create another session for the same student and exam, assert error is thrown")
	_, err = app.SessionModule.CreateSession.Execute(studentCtx, submission_dtos.CreateSessionInput{
		UserID: studentID,
		ExamID: examID,
	})
	if err == nil {
		t.Fatal("expected error when creating a duplicated session")
	}
	t.Logf("[OK] Duplicate session rejected: %v", err)

	t.Log("[STEP 10] Try to create session for non existing exam, assert error is thrown")
	_, err = app.SessionModule.CreateSession.Execute(teacherCtx, submission_dtos.CreateSessionInput{
		UserID: studentID,
		ExamID: "non-existing-exam-id",
	})
	if err == nil {
		t.Fatal("expected error when creating session with non existing exam")
	}
	t.Logf("[OK] Non existing exam session rejected: %v", err)

	t.Log("[STEP 11] Make heart beat for the session")
	// Heartbeat use case resolves active session from authenticated user in context.
	// Create a session for the authenticated teacher so heartbeat can be exercised reliably.
	teacherSession, err := app.SessionModule.CreateSession.Execute(teacherCtx, submission_dtos.CreateSessionInput{
		UserID: teacherID,
		ExamID: examID,
	})
	if err != nil {
		t.Fatalf("create teacher session for heartbeat failed: %v", err)
	}
	teacherSessionID = teacherSession.ID

	heartbeatSession, err := app.SessionModule.HeartBeatSession.Execute(teacherCtx, submission_dtos.HeartbeatSessionInput{
		SessionID: teacherSessionID,
	})
	if err != nil {
		t.Fatalf("heartbeat session failed: %v", err)
	}

	t.Log("[STEP 12] Assert session is updated successfully")
	if heartbeatSession == nil || heartbeatSession.ID == "" {
		t.Fatal("expected heartbeat session response")
	}
	if heartbeatSession.ID != teacherSessionID {
		t.Fatalf("expected heartbeat session ID=%s, got=%s", teacherSessionID, heartbeatSession.ID)
	}
	t.Logf("[OK] Heartbeat updated session successfully. sessionID=%s", heartbeatSession.ID)

	t.Log("[CLEANUP] Delete created exam, course and users")
}

func TestSubmissions(t *testing.T) {
	t.Log("[STEP 1] Initialize application container with dependencies")
	app, err := container.BuildApplicationContainer()
	if err != nil {
		t.Fatalf("failed to build application: %v", err)
	}
	t.Log("[OK] Application initialized")

	t.Log("[STEP 2] Login professor user and create context with credentials")
	teacherAccess := utils.EnsureTeacherAccess(t, app)
	teacherCtx := utils.TeacherCourseCtx(teacherAccess)
	teacherID := teacherAccess.UserData.ID
	t.Logf("[OK] Teacher login successful. teacherID=%s", teacherID)

	now := time.Now().UTC()
	enrollmentCode := fmt.Sprintf("ENR-UC-SUB-%d", now.UnixNano())
	courseCode := fmt.Sprintf("UC-SUB-%d", now.Unix()%100000)

	var courseID string
	var examID string
	var challengeID string
	testCaseIDs := make([]string, 0, 2)
	var sessionID string
	var submissionID string

	defer func() {
		if submissionID != "" {
			t.Logf("[CLEANUP] Deleting submission %s", submissionID)
			_ = app.Dependencies.SubmissionRepository.DeleteSubmission(teacherCtx, submissionID)
		}
		for _, testCaseID := range testCaseIDs {
			if testCaseID == "" {
				continue
			}
			t.Logf("[CLEANUP] Deleting test case %s", testCaseID)
			_ = app.TestCaseModule.DeleteTestCase.Execute(teacherCtx, exam_dtos.DeleteTestCaseInput{TestCaseID: testCaseID})
		}
		if challengeID != "" {
			t.Logf("[CLEANUP] Deleting challenge %s", challengeID)
			_ = app.ChallengeModule.DeleteChallenge.Execute(teacherCtx, exam_dtos.DeleteChallengeInput{ChallengeID: challengeID})
		}
		if sessionID != "" {
			t.Logf("[CLEANUP] Deleting session %s", sessionID)
			_ = app.Dependencies.SessionRepository.DeleteSession(teacherCtx, sessionID)
		}
		if examID != "" {
			t.Logf("[CLEANUP] Deleting exam %s", examID)
			_, _ = app.ExamModule.DeleteExam.Execute(teacherCtx, exam_dtos.DeleteExamInput{ExamID: examID})
		}
		if courseID != "" {
			t.Logf("[CLEANUP] Deleting course %s", courseID)
			_ = app.CourseModule.DeleteCourse.Execute(teacherCtx, course_dtos.DeleteCourseInput{CourseID: courseID})
		}
	}()

	t.Log("[STEP 3] Create a course for testing")
	createdCourse, err := app.CourseModule.CreateCourse.Execute(teacherCtx, course_dtos.CreateCourseInput{
		Name:           "UC Submission Course",
		Description:    "Course for submission use case test",
		Visibility:     string(course_entities.CourseVisibilityPublic),
		VisualIdentity: string(course_entities.CourseColourBlue),
		Code:           courseCode,
		Year:           now.Year(),
		Semester:       string(course_entities.AcademicFirstPeriod),
		EnrollmentCode: enrollmentCode,
		EnrollmentURL:  "https://example.test/enroll/" + enrollmentCode,
		TeacherID:      teacherID,
	})
	if err != nil {
		t.Fatalf("create course failed: %v", err)
	}
	if createdCourse == nil || createdCourse.ID == "" {
		t.Fatal("expected created course with ID")
	}
	courseID = createdCourse.ID
	t.Logf("[OK] Course created. courseID=%s", courseID)

	t.Log("[STEP 4] Create an exam for testing")
	startTime := now.Add(2 * time.Hour)
	endTime := startTime.Add(90 * time.Minute)
	endTimeStr := endTime.Format(time.RFC3339)

	createdExam, err := app.ExamModule.CreateExam.Execute(teacherCtx, exam_dtos.CreateExamInput{
		CourseID:             courseID,
		Title:                "UC Submission Exam",
		Description:          "Exam for submission use case test",
		Visibility:           string(exam_entities.VisibilityCourse),
		StartTime:            startTime.Format(time.RFC3339),
		EndTime:              &endTimeStr,
		AllowLateSubmissions: false,
		TimeLimit:            3600,
		TryLimit:             2,
		ProfessorID:          teacherID,
	})
	if err != nil {
		t.Fatalf("create exam failed: %v", err)
	}
	if createdExam == nil || createdExam.ID == "" {
		t.Fatal("expected created exam with ID")
	}
	examID = createdExam.ID
	t.Logf("[OK] Exam created. examID=%s", examID)

	t.Log("[STEP 5] Login student user and create context with credentials")
	studentAccess := utils.EnsureStudentAccess(t, app)
	studentCtx := utils.StudentCtx(studentAccess)
	studentID := studentAccess.UserData.ID
	t.Logf("[OK] Student access ready. studentID=%s", studentID)

	t.Log("[STEP 6] Enroll student in course")
	_, err = app.CourseModule.EnrollInCourse.Execute(studentCtx, course_dtos.EnrolledInCourseInput{
		CourseID:  courseID,
		StudentID: &studentID,
		StudentEmail: nil,
	})
	if err != nil {
		t.Fatalf("enroll student failed: %v", err)
	}
	t.Log("[OK] Student enrolled in course")

	t.Log("[STEP 6.1] Create challenge and test cases required by submission use case")
	createdChallenge, err := app.ChallengeModule.CreateChallenge.Execute(teacherCtx, exam_dtos.CreateChallengeInput{
		Title:             "UC Submission Challenge",
		Description:       "Challenge for submission use case",
		Tags:              []string{"submission", "uc"},
		Status:            string(exam_entities.ChallengeStatusPublished),
		Difficulty:        string(exam_entities.ChallengeDifficultyEasy),
		WorkerTimeLimit:   1500,
		WorkerMemoryLimit: 256,
		InputVariables: []exam_dtos.IOVariableDTO{
			{Name: "a", Type: "int", Value: "2"},
			{Name: "b", Type: "int", Value: "3"},
		},
		OutputVariable: exam_dtos.IOVariableDTO{Name: "sum", Type: "int", Value: "5"},
		Constraints:    "1 <= a,b <= 1000",
		UserID: teacherID,
	})
	if err != nil {
		t.Fatalf("create challenge failed: %v", err)
	}
	if createdChallenge == nil || createdChallenge.ID == "" {
		t.Fatal("expected created challenge with ID")
	}
	challengeID = createdChallenge.ID

	createTestCase := func(name, expected string) string {
		created, createErr := app.TestCaseModule.CreateTestCase.Execute(teacherCtx, exam_dtos.CreateTestCaseInput{
			Name: name,
			Input: []exam_dtos.IOVariableDTO{
				{Name: "a", Type: "int", Value: "2"},
				{Name: "b", Type: "int", Value: "3"},
			},
			ExpectedOutput: exam_dtos.IOVariableDTO{Name: "sum", Type: "int", Value: expected},
			IsSample:       true,
			Points:         10,
			ChallengeID:    challengeID,
		})
		if createErr != nil {
			t.Fatalf("create test case %s failed: %v", name, createErr)
		}
		if created == nil || created.ID == "" {
			t.Fatalf("expected test case id for %s", name)
		}
		return created.ID
	}

	testCaseIDs = append(testCaseIDs, createTestCase("tc_submission_1", "5"), createTestCase("tc_submission_2", "5"))
	t.Logf("[OK] Challenge and test cases created. challengeID=%s testCases=%d", challengeID, len(testCaseIDs))

	t.Log("[STEP 6.2] Create ExamItem for the challenge")
	_, err = app.ExamItemModule.CreateExamItem.Execute(teacherCtx, exam_dtos.CreateExamItemInput{
		ExamID:      examID,
		ChallengeID: challengeID,
		Order:       1,
		Points:      100,
	})

	t.Log("[STEP 7] Create Session for Student and Exam")
	createdSession, err := app.SessionModule.CreateSession.Execute(teacherCtx, submission_dtos.CreateSessionInput{
		UserID: studentID,
		ExamID: examID,
	})
	if err != nil {
		t.Fatalf("create session failed: %v", err)
	}
	if createdSession == nil || createdSession.ID == "" {
		t.Fatal("expected created session with ID")
	}
	sessionID = createdSession.ID
	t.Logf("[OK] Session created. sessionID=%s", sessionID)

	t.Log("[STEP 8] Create a submission for the session")
	createSubmissionUC := submission_usecases.NewCreateSubmissionUseCase(
		app.Dependencies.UserRepository,
		app.Dependencies.SubmissionRepository,
		app.Dependencies.SessionRepository,
		app.Dependencies.ChallengeRepository,
		app.Dependencies.TestCaseRepository,
		app.Dependencies.SubmissionResultRepository,
		app.Dependencies.PublisherPort,
	)

	createdSubmission, err := createSubmissionUC.Execute(teacherCtx, submission_dtos.CreateSubmissionInput{
		Code:        "def solve(a, b):\n    return a + b",
		Language:    string(submission_entities.LanguagePython),
		ChallengeID: challengeID,
		SessionID:   sessionID,
	})
	if err != nil {
		t.Fatalf("create submission failed: %v", err)
	}
	if createdSubmission == nil || createdSubmission.ID == "" {
		t.Fatal("expected created submission with ID")
	}
	submissionID = createdSubmission.ID
	t.Logf("[OK] Submission created. submissionID=%s", submissionID)

	t.Log("[STEP 9] Get submissions for the session and assert the created submission is returned")
	getChallengeSubmissionsUC := submission_usecases.NewGetChallengeSubmissionsUseCase(
		app.Dependencies.UserRepository,
		app.Dependencies.ChallengeRepository,
		app.Dependencies.ExamRepository,
		app.Dependencies.SubmissionRepository,
		app.Dependencies.SubmissionResultRepository,
	)

	challengeSubmissions, err := getChallengeSubmissionsUC.Execute(teacherCtx, submission_dtos.GetChallengeSubmissionsInput{
		ChallengeID: challengeID,
	})
	if err != nil {
		t.Fatalf("get challenge submissions failed: %v", err)
	}

	foundSubmission := false
	for _, output := range challengeSubmissions {
		if output == nil {
			continue
		}
		if output.Submission.ID == submissionID && output.Submission.SessionID == sessionID {
			foundSubmission = true
			break
		}
	}

	if !foundSubmission {
		t.Fatal("expected created submission to be present for the session")
	}
	t.Logf("[OK] Submission found in challenge submissions. total=%d", len(challengeSubmissions))
}
