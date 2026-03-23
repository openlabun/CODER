package usecases_test

import (
	"fmt"
	"testing"
	"time"

	course_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/course"
	exam_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	course_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"
	exam_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
)

func TestTestCasesCRUD(t *testing.T) {
	t.Log("[STEP 1] Inicializando application container y dependencias")
	app, err := buildExamApplication()
	if err != nil {
		t.Fatalf("failed to build application: %v", err)
	}
	t.Log("[OK] Application container inicializado")

	t.Log("[STEP 2] Login de profesor y obtencion de contexto autenticado")
	teacherAccess := mustLoginExamTeacher(t, app)
	teacherCtx := teacherExamCtx(teacherAccess)
	teacherID := teacherAccess.UserData.ID
	t.Logf("[OK] Login profesor exitoso. teacherID=%s", teacherID)

	now := time.Now().UTC()
	var examID string
	var challengeID string
	var testCaseID string

	defer func() {
		if testCaseID != "" {
			t.Logf("[CLEANUP] Eliminando test case pendiente %s", testCaseID)
			_ = app.TestCaseModule.DeleteTestCase.Execute(teacherCtx, exam_dtos.DeleteTestCaseInput{TestCaseID: testCaseID})
		}
		if challengeID != "" {
			t.Logf("[CLEANUP] Eliminando challenge %s", challengeID)
			_ = app.ChallengeModule.DeleteChallenge.Execute(teacherCtx, exam_dtos.DeleteChallengeInput{ChallengeID: challengeID})
		}
		if examID != "" {
			t.Logf("[CLEANUP] Eliminando examen %s", examID)
			_, _ = app.ExamModule.DeleteExam.Execute(teacherCtx, exam_dtos.DeleteExamInput{ExamID: examID})
		}
	}()

	startTime := now.Add(2 * time.Hour)
	endTime := startTime.Add(90 * time.Minute)
	endTimeStr := endTime.Format(time.RFC3339)

	t.Log("[STEP 3] Creando examen para asociar challenge")
	createdExam, err := app.ExamModule.CreateExam.Execute(teacherCtx, exam_dtos.CreateExamInput{
		CourseID:             "",
		Title:                "UC TestCase Exam",
		Description:          "Exam for test case CRUD use case test",
		Visibility:           string(exam_entities.VisibilityPrivate),
		StartTime:            startTime.Format(time.RFC3339),
		EndTime:              &endTimeStr,
		AllowLateSubmissions: false,
		TimeLimit:            5400,
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
	t.Logf("[OK] Examen creado. examID=%s", examID)

	t.Log("[STEP 4] Creando challenge para asociar test case")
	createdChallenge, err := app.ChallengeModule.CreateChallenge.Execute(teacherCtx, exam_dtos.CreateChallengeInput{
		Title:             "UC TestCase Challenge",
		Description:       "Challenge for test case CRUD use case test",
		Tags:              []string{"test-case", "crud"},
		Status:            string(exam_entities.ChallengeStatusDraft),
		Difficulty:        string(exam_entities.ChallengeDifficultyEasy),
		WorkerTimeLimit:   1500,
		WorkerMemoryLimit: 256,
		InputVariables: []exam_dtos.IOVariableDTO{
			{Name: "a", Type: "int", Value: "2"},
			{Name: "b", Type: "int", Value: "3"},
		},
		OutputVariable: exam_dtos.IOVariableDTO{Name: "sum", Type: "int", Value: "5"},
		Constraints:    "1 <= a,b <= 1000",
		ExamID:          examID,
	})
	if err != nil {
		t.Fatalf("create challenge failed: %v", err)
	}
	if createdChallenge == nil || createdChallenge.ID == "" {
		t.Fatal("expected created challenge with ID")
	}
	challengeID = createdChallenge.ID
	t.Logf("[OK] Challenge creado. challengeID=%s", challengeID)

	t.Log("[STEP 5][TEST 1] Creando test case")
	createdTestCase, err := app.TestCaseModule.CreateTestCase.Execute(teacherCtx, exam_dtos.CreateTestCaseInput{
		Name: "tc_add_1",
		Input: []exam_dtos.IOVariableDTO{
			{Name: "a", Type: "int", Value: "10"},
			{Name: "b", Type: "int", Value: "15"},
		},
		ExpectedOutput: exam_dtos.IOVariableDTO{Name: "sum", Type: "int", Value: "25"},
		IsSample:       true,
		Points:         10,
		ChallengeID:    challengeID,
	})
	if err != nil {
		t.Fatalf("create test case failed: %v", err)
	}
	if createdTestCase == nil || createdTestCase.ID == "" {
		t.Fatal("expected created test case with ID")
	}
	testCaseID = createdTestCase.ID
	t.Logf("[OK] Test case creado. testCaseID=%s", testCaseID)

	t.Log("[STEP 6][TEST 2] Obteniendo test cases por challenge")
	testCasesByChallenge, err := app.TestCaseModule.GetTestCasesByChallenge.Execute(teacherCtx, exam_dtos.GetTestCasesByChallengeInput{ChallengeID: challengeID})
	if err != nil {
		t.Fatalf("get test cases by challenge failed: %v", err)
	}
	foundCreated := false
	for _, tc := range testCasesByChallenge {
		if tc != nil && tc.ID == testCaseID {
			foundCreated = true
			break
		}
	}
	if !foundCreated {
		t.Fatal("expected created test case in challenge list")
	}
	t.Logf("[OK] Test case encontrado en listado. totalTestCases=%d", len(testCasesByChallenge))

	t.Log("[STEP 7][TEST 3] Actualizando test case")
	updatedName := "tc_add_1_updated"
	updatedPoints := 25
	updatedIsSample := false
	updatedTestCase, err := app.TestCaseModule.UpdateTestCase.Execute(teacherCtx, exam_dtos.UpdateTestCaseInput{
		ID:       testCaseID,
		Name:     &updatedName,
		Points:   &updatedPoints,
		IsSample: &updatedIsSample,
	})
	if err != nil {
		t.Fatalf("update test case failed: %v", err)
	}
	if updatedTestCase == nil || updatedTestCase.Name != updatedName || updatedTestCase.Points != updatedPoints || updatedTestCase.IsSample != updatedIsSample {
		t.Fatal("expected updated test case values")
	}
	t.Logf("[OK] Test case actualizado. name=%q points=%d isSample=%v", updatedTestCase.Name, updatedTestCase.Points, updatedTestCase.IsSample)

	t.Log("[STEP 8][TEST 4] Obteniendo test cases por challenge y validando update")
	reloadedList, err := app.TestCaseModule.GetTestCasesByChallenge.Execute(teacherCtx, exam_dtos.GetTestCasesByChallengeInput{ChallengeID: challengeID})
	if err != nil {
		t.Fatalf("get test cases after update failed: %v", err)
	}
	validatedUpdate := false
	for _, tc := range reloadedList {
		if tc != nil && tc.ID == testCaseID {
			if tc.Name != updatedName || tc.Points != updatedPoints || tc.IsSample != updatedIsSample {
				t.Fatalf("expected persisted update for test case %s", testCaseID)
			}
			validatedUpdate = true
			break
		}
	}
	if !validatedUpdate {
		t.Fatal("expected updated test case in reloaded list")
	}
	t.Log("[OK] Persistencia de update validada")

	t.Logf("[STEP 9][TEST 5] Eliminando test case testCaseID=%s", testCaseID)
	if err := app.TestCaseModule.DeleteTestCase.Execute(teacherCtx, exam_dtos.DeleteTestCaseInput{TestCaseID: testCaseID}); err != nil {
		t.Fatalf("delete test case failed: %v", err)
	}
	testCaseID = ""

	afterDelete, err := app.TestCaseModule.GetTestCasesByChallenge.Execute(teacherCtx, exam_dtos.GetTestCasesByChallengeInput{ChallengeID: challengeID})
	if err != nil {
		t.Fatalf("get test cases after delete failed: %v", err)
	}
	for _, tc := range afterDelete {
		if tc != nil && tc.ID == createdTestCase.ID {
			t.Fatal("expected deleted test case to be absent from challenge list")
		}
	}
	t.Log("[OK] Test case eliminado y ausencia validada")

	t.Log("[STEP 10][CLEANUP] Limpieza automatica de challenge y examen preparada")
}

func TestTestCasesFromStudentView(t *testing.T) {
	t.Log("[STEP 1] Inicializando application container y dependencias")
	app, err := buildExamApplication()
	if err != nil {
		t.Fatalf("failed to build application: %v", err)
	}
	t.Log("[OK] Application container inicializado")

	t.Log("[STEP 2] Login de profesor y contexto autenticado")
	teacherAccess := mustLoginExamTeacher(t, app)
	teacherCtx := teacherExamCtx(teacherAccess)
	teacherID := teacherAccess.UserData.ID
	t.Logf("[OK] Login profesor exitoso. teacherID=%s", teacherID)

	now := time.Now().UTC()
	enrollmentCode := fmt.Sprintf("ENR-UC-TC-STU-%d", now.UnixNano())
	courseCode := fmt.Sprintf("UC-TC-STU-%d", now.Unix()%100000)

	var courseID string
	var examID string
	var challengeID string
	testCaseIDs := make([]string, 0, 3)

	defer func() {
		for _, id := range testCaseIDs {
			if id == "" {
				continue
			}
			t.Logf("[CLEANUP] Eliminando test case %s", id)
			_ = app.TestCaseModule.DeleteTestCase.Execute(teacherCtx, exam_dtos.DeleteTestCaseInput{TestCaseID: id})
		}
		if challengeID != "" {
			t.Logf("[CLEANUP] Eliminando challenge %s", challengeID)
			_ = app.ChallengeModule.DeleteChallenge.Execute(teacherCtx, exam_dtos.DeleteChallengeInput{ChallengeID: challengeID})
		}
		if examID != "" {
			t.Logf("[CLEANUP] Eliminando examen %s", examID)
			_, _ = app.ExamModule.DeleteExam.Execute(teacherCtx, exam_dtos.DeleteExamInput{ExamID: examID})
		}
		if courseID != "" {
			t.Logf("[CLEANUP] Eliminando curso %s", courseID)
			_ = app.CourseModule.DeleteCourse.Execute(teacherCtx, course_dtos.DeleteCourseInput{CourseID: courseID})
		}
	}()

	t.Logf("[STEP 3] Creando curso para asociar examen code=%s", courseCode)
	createdCourse, err := app.CourseModule.CreateCourse.Execute(teacherCtx, course_dtos.CreateCourseInput{
		Name:           "UC TestCase Student View Course",
		Description:    "Course for test case student visibility",
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
	t.Logf("[OK] Curso creado. courseID=%s", courseID)

	startTime := time.Now().UTC().Add(2 * time.Hour)
	endTime := startTime.Add(90 * time.Minute)
	endTimeStr := endTime.Format(time.RFC3339)

	t.Log("[STEP 4] Creando examen asociado al curso")
	createdExam, err := app.ExamModule.CreateExam.Execute(teacherCtx, exam_dtos.CreateExamInput{
		CourseID:             courseID,
		Title:                "UC TestCase Student View Exam",
		Description:          "Exam for student visibility restrictions",
		Visibility:           string(exam_entities.VisibilityCourse),
		StartTime:            startTime.Format(time.RFC3339),
		EndTime:              &endTimeStr,
		AllowLateSubmissions: false,
		TimeLimit:            5400,
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
	t.Logf("[OK] Examen creado. examID=%s", examID)

	t.Log("[STEP 5] Creando challenge para asociar test cases")
	createdChallenge, err := app.ChallengeModule.CreateChallenge.Execute(teacherCtx, exam_dtos.CreateChallengeInput{
		Title:             "UC TestCase Visibility Challenge",
		Description:       "Challenge for test case visibility test",
		Tags:              []string{"visibility", "student"},
		Status:            string(exam_entities.ChallengeStatusPublished),
		Difficulty:        string(exam_entities.ChallengeDifficultyEasy),
		WorkerTimeLimit:   1500,
		WorkerMemoryLimit: 256,
		InputVariables: []exam_dtos.IOVariableDTO{
			{Name: "x", Type: "int", Value: "1"},
		},
		OutputVariable: exam_dtos.IOVariableDTO{Name: "y", Type: "int", Value: "1"},
		Constraints:    "1 <= x <= 10",
		ExamID:          examID,
	})
	if err != nil {
		t.Fatalf("create challenge failed: %v", err)
	}
	if createdChallenge == nil || createdChallenge.ID == "" {
		t.Fatal("expected created challenge with ID")
	}
	challengeID = createdChallenge.ID
	t.Logf("[OK] Challenge creado. challengeID=%s", challengeID)

	t.Log("[STEP 6][TEST 1] Creando 3 test cases: 2 publicos (IsSample=true) y 1 privado (IsSample=false)")
	createTestCase := func(name string, isSample bool) string {
		created, createErr := app.TestCaseModule.CreateTestCase.Execute(teacherCtx, exam_dtos.CreateTestCaseInput{
			Name: name,
			Input: []exam_dtos.IOVariableDTO{
				{Name: "x", Type: "int", Value: "2"},
			},
			ExpectedOutput: exam_dtos.IOVariableDTO{Name: "y", Type: "int", Value: "2"},
			IsSample:       isSample,
			Points:         10,
			ChallengeID:    challengeID,
		})
		if createErr != nil {
			t.Fatalf("create test case %s failed: %v", name, createErr)
		}
		if created == nil || created.ID == "" {
			t.Fatalf("expected created test case with ID for %s", name)
		}
		return created.ID
	}

	publicTC1 := createTestCase("tc_public_1", true)
	publicTC2 := createTestCase("tc_public_2", true)
	privateTC := createTestCase("tc_private_1", false)
	testCaseIDs = append(testCaseIDs, publicTC1, publicTC2, privateTC)
	t.Logf("[OK] Test cases creados. public1=%s public2=%s private=%s", publicTC1, publicTC2, privateTC)

	t.Log("[STEP 7] Login/registro de estudiante y contexto autenticado")
	studentAccess := ensureExamStudentAccess(t, app)
	studentCtx := studentExamCtx(studentAccess)
	t.Logf("[OK] Estudiante listo. studentID=%s", studentAccess.UserData.ID)

	t.Log("[STEP 8] Matriculando estudiante en el curso")
	_, err = app.CourseModule.EnrollInCourse.Execute(studentCtx, course_dtos.EnrolledInCourseInput{
		CourseID:  courseID,
		StudentID: studentAccess.UserData.ID,
	})
	if err != nil {
		t.Fatalf("enroll student failed: %v", err)
	}
	t.Log("[OK] Estudiante matriculado")

	t.Log("[STEP 9][TEST 2] Obteniendo test cases desde vista de estudiante y validando restricciones")
	studentTestCases, err := app.TestCaseModule.GetTestCasesByChallenge.Execute(studentCtx, exam_dtos.GetTestCasesByChallengeInput{ChallengeID: challengeID})
	if err != nil {
		t.Fatalf("get test cases by challenge as student failed: %v", err)
	}

	foundPublic1 := false
	foundPublic2 := false
	foundPrivate := false
	for _, tc := range studentTestCases {
		if tc == nil {
			continue
		}
		if tc.ID == publicTC1 {
			foundPublic1 = true
		}
		if tc.ID == publicTC2 {
			foundPublic2 = true
		}
		if tc.ID == privateTC {
			foundPrivate = true
		}
	}

	if !foundPublic1 || !foundPublic2 {
		t.Fatal("expected both public test cases in student view, got only public1:", foundPublic1, "public2:", foundPublic2)
	}
	if foundPrivate {
		t.Fatal("did not expect private test case in student view")
	}
	t.Logf("[OK] Restricciones de vista estudiante validadas. visibleTestCases=%d", len(studentTestCases))

	t.Log("[STEP 10][CLEANUP] Limpieza automatica de test cases/challenge/exam/course preparada")
}