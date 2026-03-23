package usecases_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	course_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"
	exam_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"

	container "github.com/openlabun/CODER/apps/api_v2/internal/application/container"
	course_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/course"
	exam_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	user_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/user"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
)

func TestChallengeCRUD(t *testing.T) {
	t.Log("[STEP 1] Inicializando application container con dependencias")
	app, err := buildExamApplication()
	if err != nil {
		t.Fatalf("failed to build application: %v", err)
	}
	t.Log("[OK] Application container inicializado")

	t.Log("[STEP 2] Login de profesor y construccion de contexto autenticado")
	teacherAccess := mustLoginExamTeacher(t, app)
	teacherCtx := teacherExamCtx(teacherAccess)
	teacherID := teacherAccess.UserData.ID
	t.Logf("[OK] Login profesor exitoso. teacherID=%s", teacherID)

	now := time.Now().UTC()
	enrollmentCode := fmt.Sprintf("ENR-UC-CH-%d", now.UnixNano())
	courseCode := fmt.Sprintf("UC-CH-%d", now.Unix()%100000)

	t.Logf("[STEP 3] Creando curso de pruebas code=%s", courseCode)
	createdCourse, err := app.CourseModule.CreateCourse.Execute(teacherCtx, course_dtos.CreateCourseInput{
		Name:           "UC Challenge Course",
		Description:    "Course for challenge CRUD use case test",
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
	courseID := createdCourse.ID
	t.Logf("[OK] Curso creado. courseID=%s", courseID)

	var examID string
	var challengeID string
	defer func() {
		if challengeID != "" {
			t.Logf("[CLEANUP] Eliminando challenge pendiente %s", challengeID)
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

	startTime := time.Now().UTC().Add(2 * time.Hour)
	endTime := startTime.Add(90 * time.Minute)
	endTimeStr := endTime.Format(time.RFC3339)

	t.Log("[STEP 4] Creando examen de pruebas")
	createdExam, err := app.ExamModule.CreateExam.Execute(teacherCtx, exam_dtos.CreateExamInput{
		CourseID:             courseID,
		Title:                "UC Challenge Exam",
		Description:          "Exam for challenge CRUD use case test",
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

	t.Log("[STEP 5] Creando challenge de pruebas")
	createdChallenge, err := app.ChallengeModule.CreateChallenge.Execute(teacherCtx, exam_dtos.CreateChallengeInput{
		Title:             "UC Challenge",
		Description:       "Challenge created by use case test",
		Tags:              []string{"algorithms", "uc"},
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

	t.Log("[STEP 6] Actualizando challenge y validando cambios")
	updatedTitle := "UC Challenge Updated"
	updatedDescription := "Challenge updated by use case test"
	updatedTags := []string{"algorithms", "updated"}
	updatedTimeLimit := 2000
	updatedMemoryLimit := 512
	updatedConstraints := "1 <= a,b <= 2000"
	updatedChallenge, err := app.ChallengeModule.UpdateChallenge.Execute(teacherCtx, exam_dtos.UpdateChallengeInput{
		ChallengeID:       challengeID,
		Title:             &updatedTitle,
		Description:       &updatedDescription,
		Tags:              &updatedTags,
		WorkerTimeLimit:   &updatedTimeLimit,
		WorkerMemoryLimit: &updatedMemoryLimit,
		Constraints:       &updatedConstraints,
	})
	if err != nil {
		t.Fatalf("update challenge failed: %v", err)
	}
	if updatedChallenge == nil || updatedChallenge.Title != updatedTitle {
		t.Fatal("expected updated challenge title")
	}
	t.Logf("[OK] Challenge actualizado. title=%q", updatedChallenge.Title)

	t.Log("[STEP 7] Publicando challenge y validando estado")
	publishedChallenge, err := app.ChallengeModule.PublishChallenge.Execute(teacherCtx, exam_dtos.PublishChallengeInput{ChallengeID: challengeID})
	if err != nil {
		t.Fatalf("publish challenge failed: %v", err)
	}
	if publishedChallenge == nil || publishedChallenge.Status != exam_entities.ChallengeStatusPublished {
		t.Fatal("expected challenge status published")
	}
	t.Logf("[OK] Challenge publicado. status=%s", publishedChallenge.Status)

	t.Logf("[STEP 8] Listando challenges del examen examID=%s", examID)
	challengesByExam, err := app.ChallengeModule.GetChallengesByExam.Execute(teacherCtx, exam_dtos.GetChallengesByExamInput{ExamID: examID})
	if err != nil {
		t.Fatalf("get challenges by exam failed: %v", err)
	}
	found := false
	for _, ch := range challengesByExam {
		if ch != nil && ch.ID == challengeID {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected created challenge in exam challenge list")
	}
	t.Logf("[OK] Challenge encontrado en listado. totalChallenges=%d", len(challengesByExam))

	t.Log("[STEP 9] Archivando challenge y validando estado")
	archivedChallenge, err := app.ChallengeModule.ArchiveChallenge.Execute(teacherCtx, exam_dtos.ArchiveChallengeInput{ChallengeID: challengeID})
	if err != nil {
		t.Fatalf("archive challenge failed: %v", err)
	}
	if archivedChallenge == nil || archivedChallenge.Status != exam_entities.ChallengeStatusArchived {
		t.Fatal("expected challenge status archived")
	}
	t.Logf("[OK] Challenge archivado. status=%s", archivedChallenge.Status)

	t.Logf("[STEP 10] Eliminando challenge challengeID=%s", challengeID)
	if err := app.ChallengeModule.DeleteChallenge.Execute(teacherCtx, exam_dtos.DeleteChallengeInput{ChallengeID: challengeID}); err != nil {
		t.Fatalf("delete challenge failed: %v", err)
	}
	challengeID = ""

	challengesAfterDelete, err := app.ChallengeModule.GetChallengesByExam.Execute(teacherCtx, exam_dtos.GetChallengesByExamInput{ExamID: examID})
	if err != nil {
		t.Fatalf("get challenges after delete failed: %v", err)
	}
	for _, ch := range challengesAfterDelete {
		if ch != nil && ch.ID == updatedChallenge.ID {
			t.Fatal("expected deleted challenge to be absent from exam challenge list")
		}
	}
	t.Log("[OK] Challenge eliminado y ausencia validada")

	t.Log("[STEP 11] Teardown automatico preparado (exam + course)")

}

func TestChallengeFromStudentView(t *testing.T) {
	t.Log("[STEP 1] Inicializando application container con dependencias")
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
	enrollmentCode := fmt.Sprintf("ENR-UC-CH-STU-%d", now.UnixNano())
	courseCode := fmt.Sprintf("UC-CH-STU-%d", now.Unix()%100000)

	t.Logf("[STEP 3] Creando curso de pruebas code=%s", courseCode)
	createdCourse, err := app.CourseModule.CreateCourse.Execute(teacherCtx, course_dtos.CreateCourseInput{
		Name:           "UC Challenge Student View Course",
		Description:    "Course for student-view challenge restrictions",
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
	courseID := createdCourse.ID
	t.Logf("[OK] Curso creado. courseID=%s", courseID)

	var examID string
	publishedChallengeID := ""
	archivedChallengeID := ""
	draftChallengeID := ""
	defer func() {
		if publishedChallengeID != "" {
			t.Logf("[CLEANUP] Eliminando challenge publicado %s", publishedChallengeID)
			_ = app.ChallengeModule.DeleteChallenge.Execute(teacherCtx, exam_dtos.DeleteChallengeInput{ChallengeID: publishedChallengeID})
		}
		if archivedChallengeID != "" {
			t.Logf("[CLEANUP] Eliminando challenge archivado %s", archivedChallengeID)
			_ = app.ChallengeModule.DeleteChallenge.Execute(teacherCtx, exam_dtos.DeleteChallengeInput{ChallengeID: archivedChallengeID})
		}
		if draftChallengeID != "" {
			t.Logf("[CLEANUP] Eliminando challenge draft %s", draftChallengeID)
			_ = app.ChallengeModule.DeleteChallenge.Execute(teacherCtx, exam_dtos.DeleteChallengeInput{ChallengeID: draftChallengeID})
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

	startTime := time.Now().UTC().Add(2 * time.Hour)
	endTime := startTime.Add(90 * time.Minute)
	endTimeStr := endTime.Format(time.RFC3339)

	t.Log("[STEP 4] Creando examen accesible para estudiantes")
	createdExam, err := app.ExamModule.CreateExam.Execute(teacherCtx, exam_dtos.CreateExamInput{
		CourseID:             courseID,
		Title:                "UC Challenge Student View Exam",
		Description:          "Exam for student-view challenge restrictions",
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

	t.Log("[STEP 5] Creando 3 challenges (draft inicial)")
	createDraftChallenge := func(title string) string {
		created, createErr := app.ChallengeModule.CreateChallenge.Execute(teacherCtx, exam_dtos.CreateChallengeInput{
			Title:             title,
			Description:       "Challenge for student view restrictions",
			Tags:              []string{"student-view"},
			Status:            string(exam_entities.ChallengeStatusDraft),
			Difficulty:        string(exam_entities.ChallengeDifficultyEasy),
			WorkerTimeLimit:   1500,
			WorkerMemoryLimit: 256,
			InputVariables: []exam_dtos.IOVariableDTO{
				{Name: "a", Type: "int", Value: "1"},
			},
			OutputVariable: exam_dtos.IOVariableDTO{Name: "out", Type: "int", Value: "1"},
			Constraints:    "1 <= a <= 10",
			ExamID:          examID,
		})
		if createErr != nil {
			t.Fatalf("create challenge %s failed: %v", title, createErr)
		}
		if created == nil || created.ID == "" {
			t.Fatalf("expected created challenge with ID for %s", title)
		}
		return created.ID
	}

	publishedChallengeID = createDraftChallenge("Challenge Published")
	archivedChallengeID = createDraftChallenge("Challenge Archived")
	draftChallengeID = createDraftChallenge("Challenge Draft")
	t.Logf("[OK] Challenges creados. publishedCandidate=%s archivedCandidate=%s draft=%s", publishedChallengeID, archivedChallengeID, draftChallengeID)

	t.Log("[STEP 6] Publicando y archivando challenges para dejar estados finales")
	publishedChallenge, err := app.ChallengeModule.PublishChallenge.Execute(teacherCtx, exam_dtos.PublishChallengeInput{ChallengeID: publishedChallengeID})
	if err != nil {
		t.Fatalf("publish published-candidate failed: %v", err)
	}
	if publishedChallenge == nil || publishedChallenge.Status != exam_entities.ChallengeStatusPublished {
		t.Fatal("expected published challenge status")
	}

	archivedAfterPublish, err := app.ChallengeModule.PublishChallenge.Execute(teacherCtx, exam_dtos.PublishChallengeInput{ChallengeID: archivedChallengeID})
	if err != nil {
		t.Fatalf("publish archived-candidate failed: %v", err)
	}
	if archivedAfterPublish == nil || archivedAfterPublish.Status != exam_entities.ChallengeStatusPublished {
		t.Fatal("expected archived candidate to be published before archive")
	}

	archivedChallenge, err := app.ChallengeModule.ArchiveChallenge.Execute(teacherCtx, exam_dtos.ArchiveChallengeInput{ChallengeID: archivedChallengeID})
	if err != nil {
		t.Fatalf("archive challenge failed: %v", err)
	}
	if archivedChallenge == nil || archivedChallenge.Status != exam_entities.ChallengeStatusArchived {
		t.Fatal("expected archived challenge status")
	}
	t.Logf("[OK] Estados finales: published=%s archived=%s draft=%s", publishedChallenge.Status, archivedChallenge.Status, exam_entities.ChallengeStatusDraft)

	t.Log("[STEP 7] Login/registro de estudiante y construccion de contexto autenticado")
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

	t.Log("[STEP 9] Obteniendo challenges desde vista de estudiante y validando restricciones")
	studentChallenges, err := app.ChallengeModule.GetChallengesByExam.Execute(studentCtx, exam_dtos.GetChallengesByExamInput{ExamID: examID})
	if err != nil {
		t.Fatalf("get challenges by exam as student failed: %v", err)
	}

	foundPublished := false
	foundArchived := false
	foundDraft := false
	for _, ch := range studentChallenges {
		if ch == nil {
			continue
		}
		if ch.ID == publishedChallengeID {
			foundPublished = true
		}
		if ch.ID == archivedChallengeID {
			foundArchived = true
		}
		if ch.ID == draftChallengeID {
			foundDraft = true
		}
		if ch.Status != exam_entities.ChallengeStatusPublished {
			t.Fatalf("expected only published challenges in student view, got status=%s", ch.Status)
		}
	}

	if !foundPublished {
		t.Fatal("expected published challenge in student view")
	}
	if foundArchived {
		t.Fatal("did not expect archived challenge in student view")
	}
	if foundDraft {
		t.Fatal("did not expect draft challenge in student view")
	}
	t.Logf("[OK] Restricciones de vista estudiante validadas. visibleChallenges=%d", len(studentChallenges))

	t.Log("[STEP 10] Teardown automatico preparado (challenges + exam + course)")
}

func ensureExamStudentAccess(t *testing.T, app *container.Application) *user_dtos.UserAccess {
	t.Helper()
	t.Log("[AUTH] Intentando login de estudiante")

	email := "stud@test.com"
	password := "Testing123!"

	access, err := app.Dependencies.LoginService.LoginUser(email, password)
	if err == nil && access != nil && access.UserData != nil && access.UserData.ID != "" && access.Token != nil && access.Token.AccessToken != "" {
		t.Logf("[AUTH][OK] Login estudiante exitoso. studentID=%s", access.UserData.ID)
		return access
	}

	t.Log("[AUTH] Estudiante no existe o login fallo, registrando usuario")
	registered, registerErr := app.Dependencies.RegisterService.RegisterUserDirect(email, password, "Student Test")
	if registerErr != nil {
		t.Fatalf("register student failed: %v", registerErr)
	}
	if registered == nil || registered.UserData == nil || registered.UserData.ID == "" || registered.Token == nil || registered.Token.AccessToken == "" {
		t.Fatal("expected registered student with ID and access token")
	}
	t.Logf("[AUTH][OK] Registro de estudiante exitoso. studentID=%s", registered.UserData.ID)

	return registered
}

func studentExamCtx(studentAccess *user_dtos.UserAccess) context.Context {
	ctx := services.WithAccessToken(context.Background(), studentAccess.Token.AccessToken)
	ctx = services.WithUserEmail(ctx, studentAccess.UserData.Email)
	return ctx
}