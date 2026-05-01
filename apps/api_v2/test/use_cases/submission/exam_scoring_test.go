package usecases_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	exam_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	submission_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"
	exam_consts "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/exam"
	submission_consts "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/submission"
	test "github.com/openlabun/CODER/apps/api_v2/test"
	utils "github.com/openlabun/CODER/apps/api_v2/test/use_cases"
)

func TestExamScoring(t *testing.T) {
	process := test.StartTestWithApp(t, "Exam Scoring")
	teacherEmail := "test@test.com"
	studentEmail := "stud@test.com"
	password := "Password123!"

	var teacherCtx = context.Background()
	var studentCtx = context.Background()
	var studentID string
	var examID string
	var firstChallengeID string
	var secondChallengeID string
	var firstTestCaseOneID string
	var firstTestCaseTwoID string
	var firstTestCaseThreeID string
	var secondTestCaseOneID string
	var secondTestCaseTwoID string
	var firstExamItemID string
	var secondExamItemID string
	var sessionID string
	var submissionID string

	defer func() {
		if sessionID != "" {
			t.Logf("[CLEANUP] Cerrando sesión %s", sessionID)
			_, _ = process.Application.SessionModule.CloseSession.Execute(teacherCtx, submission_dtos.CloseSessionInput{SessionID: sessionID})
		}
		if firstExamItemID != "" {
			t.Logf("[CLEANUP] Eliminando exam item 1 %s", firstExamItemID)
			_ = process.Application.ExamItemModule.DeleteExamItem.Execute(teacherCtx, exam_dtos.DeleteExamItemInput{ID: firstExamItemID})
		}
		if secondExamItemID != "" {
			t.Logf("[CLEANUP] Eliminando exam item 2 %s", secondExamItemID)
			_ = process.Application.ExamItemModule.DeleteExamItem.Execute(teacherCtx, exam_dtos.DeleteExamItemInput{ID: secondExamItemID})
		}
		if firstTestCaseThreeID != "" {
			t.Logf("[CLEANUP] Eliminando test case 1.3 %s", firstTestCaseThreeID)
			_ = process.Application.TestCaseModule.DeleteTestCase.Execute(teacherCtx, exam_dtos.DeleteTestCaseInput{TestCaseID: firstTestCaseThreeID})
		}
		if firstTestCaseTwoID != "" {
			t.Logf("[CLEANUP] Eliminando test case 1.2 %s", firstTestCaseTwoID)
			_ = process.Application.TestCaseModule.DeleteTestCase.Execute(teacherCtx, exam_dtos.DeleteTestCaseInput{TestCaseID: firstTestCaseTwoID})
		}
		if firstTestCaseOneID != "" {
			t.Logf("[CLEANUP] Eliminando test case 1.1 %s", firstTestCaseOneID)
			_ = process.Application.TestCaseModule.DeleteTestCase.Execute(teacherCtx, exam_dtos.DeleteTestCaseInput{TestCaseID: firstTestCaseOneID})
		}
		if secondTestCaseTwoID != "" {
			t.Logf("[CLEANUP] Eliminando test case 2.2 %s", secondTestCaseTwoID)
			_ = process.Application.TestCaseModule.DeleteTestCase.Execute(teacherCtx, exam_dtos.DeleteTestCaseInput{TestCaseID: secondTestCaseTwoID})
		}
		if secondTestCaseOneID != "" {
			t.Logf("[CLEANUP] Eliminando test case 2.1 %s", secondTestCaseOneID)
			_ = process.Application.TestCaseModule.DeleteTestCase.Execute(teacherCtx, exam_dtos.DeleteTestCaseInput{TestCaseID: secondTestCaseOneID})
		}
		if firstChallengeID != "" {
			t.Logf("[CLEANUP] Eliminando challenge 1 %s", firstChallengeID)
			_ = process.Application.ChallengeModule.DeleteChallenge.Execute(teacherCtx, exam_dtos.DeleteChallengeInput{ChallengeID: firstChallengeID})
		}
		if secondChallengeID != "" {
			t.Logf("[CLEANUP] Eliminando challenge 2 %s", secondChallengeID)
			_ = process.Application.ChallengeModule.DeleteChallenge.Execute(teacherCtx, exam_dtos.DeleteChallengeInput{ChallengeID: secondChallengeID})
		}
		if examID != "" {
			t.Logf("[CLEANUP] Eliminando examen %s", examID)
			_, _ = process.Application.ExamModule.DeleteExam.Execute(teacherCtx, exam_dtos.DeleteExamInput{ExamID: examID})
		}
	}()

	process.StartStep("Iniciar sesión con usuario de docente")
	teacherAccess := utils.EnsureAuthUserAccess(t, process.Application, teacherEmail, password, "Teacher Test")
	teacherCtx = utils.BuildUserCtx(teacherAccess)
	process.EndStep()

	process.StartStep("Crear examen público (visibilidad public y sin curso)")
	now := time.Now().UTC()
	createdExam, err := process.Application.ExamModule.CreateExam.Execute(teacherCtx, exam_dtos.CreateExamInput{
		CourseID:             nil,
		Title:                "Exam Scoring Exam",
		Description:          "Examen para validar puntaje del examen",
		Visibility:           string(exam_consts.VisibilityPublic),
		StartTime:            now.Add(2 * time.Hour).Format(time.RFC3339),
		EndTime:              nil,
		AllowLateSubmissions: true,
		TimeLimit:            120,
		TryLimit:             3,
		ProfessorID:          teacherAccess.UserData.ID,
	})
	if err != nil {
		process.Fail("create exam", err)
	}
	examID = createdExam.ID
	process.EndStep()

	process.StartStep("Crear dos retos")
	firstChallenge, err := process.Application.ChallengeModule.CreateChallenge.Execute(teacherCtx, exam_dtos.CreateChallengeInput{
		Title:             "Exam Scoring Challenge 1",
		Description:       "Primer reto para validar score del examen",
		Tags:              []string{"exam", "scoring", "one"},
		Status:            string(exam_consts.ChallengeStatusPublished),
		Difficulty:        string(exam_consts.ChallengeDifficultyEasy),
		WorkerTimeLimit:   1200,
		WorkerMemoryLimit: 256,
		CodeTemplates: []exam_dtos.CodeTemplateDTO{{
			Language: "python",
			Template: "def solve():\n    pass",
		}},
		InputVariables: []exam_dtos.IOVariableDTO{{
			Name:  "n",
			Type:  string(exam_consts.VariableFormatInt),
			Value: "2",
		}},
		OutputVariable: exam_dtos.IOVariableDTO{Name: "out", Type: string(exam_consts.VariableFormatInt), Value: "4"},
		Constraints:    "1 <= n <= 1000",
	})
	if err != nil {
		process.Fail("create first challenge", err)
	}
	firstChallengeID = firstChallenge.ID

	secondChallenge, err := process.Application.ChallengeModule.CreateChallenge.Execute(teacherCtx, exam_dtos.CreateChallengeInput{
		Title:             "Exam Scoring Challenge 2",
		Description:       "Segundo reto para validar score del examen",
		Tags:              []string{"exam", "scoring", "two"},
		Status:            string(exam_consts.ChallengeStatusPublished),
		Difficulty:        string(exam_consts.ChallengeDifficultyEasy),
		WorkerTimeLimit:   1200,
		WorkerMemoryLimit: 256,
		CodeTemplates: []exam_dtos.CodeTemplateDTO{{
			Language: "python",
			Template: "def solve():\n    pass",
		}},
		InputVariables: []exam_dtos.IOVariableDTO{{
			Name:  "n",
			Type:  string(exam_consts.VariableFormatInt),
			Value: "5",
		}},
		OutputVariable: exam_dtos.IOVariableDTO{Name: "out", Type: string(exam_consts.VariableFormatInt), Value: "10"},
		Constraints:    "1 <= n <= 1000",
	})
	if err != nil {
		process.Fail("create second challenge", err)
	}
	secondChallengeID = secondChallenge.ID
	process.EndStep()

	process.StartStep("Crear 2 casos de prueba con valor de 3 puntos para el primer reto")
	firstTestCaseOne, err := process.Application.TestCaseModule.CreateTestCase.Execute(teacherCtx, exam_dtos.CreateTestCaseInput{
		Name:           "exam_score_case_1_1",
		Input:          []exam_dtos.IOVariableDTO{{Name: "n", Type: string(exam_consts.VariableFormatInt), Value: "2"}},
		ExpectedOutput: exam_dtos.IOVariableDTO{Name: "out", Type: string(exam_consts.VariableFormatInt), Value: "4"},
		IsSample:       false,
		Points:         3,
		ChallengeID:    firstChallengeID,
	})
	if err != nil {
		process.Fail("create first test case 1", err)
	}
	firstTestCaseOneID = firstTestCaseOne.ID

	firstTestCaseTwo, err := process.Application.TestCaseModule.CreateTestCase.Execute(teacherCtx, exam_dtos.CreateTestCaseInput{
		Name:           "exam_score_case_1_2",
		Input:          []exam_dtos.IOVariableDTO{{Name: "n", Type: string(exam_consts.VariableFormatInt), Value: "5"}},
		ExpectedOutput: exam_dtos.IOVariableDTO{Name: "out", Type: string(exam_consts.VariableFormatInt), Value: "10"},
		IsSample:       false,
		Points:         3,
		ChallengeID:    firstChallengeID,
	})
	if err != nil {
		process.Fail("create first test case 2", err)
	}
	firstTestCaseTwoID = firstTestCaseTwo.ID
	process.EndStep()

	process.StartStep("Crear un caso de prueba con valor de 6 puntos para el primer reto (debe ser imposible de cumplir)")
	firstTestCaseThree, err := process.Application.TestCaseModule.CreateTestCase.Execute(teacherCtx, exam_dtos.CreateTestCaseInput{
		Name:           "exam_score_case_1_impossible",
		Input:          []exam_dtos.IOVariableDTO{{Name: "n", Type: string(exam_consts.VariableFormatInt), Value: "7"}},
		ExpectedOutput: exam_dtos.IOVariableDTO{Name: "out", Type: string(exam_consts.VariableFormatInt), Value: "999"},
		IsSample:       false,
		Points:         6,
		ChallengeID:    firstChallengeID,
	})
	if err != nil {
		process.Fail("create first impossible test case", err)
	}
	firstTestCaseThreeID = firstTestCaseThree.ID
	process.EndStep()

	process.StartStep("Crear un punto de examen con valor de 30")
	firstExamItem, err := process.Application.ExamItemModule.CreateExamItem.Execute(teacherCtx, exam_dtos.CreateExamItemInput{
		ExamID:      examID,
		ChallengeID: firstChallengeID,
		Order:       1,
		Points:      30,
	})
	if err != nil {
		process.Fail("create first exam item", err)
	}
	firstExamItemID = firstExamItem.ID
	process.EndStep()

	process.StartStep("Crear 2 casos de prueba con valor de 5 puntos para el segundo reto")
	secondTestCaseOne, err := process.Application.TestCaseModule.CreateTestCase.Execute(teacherCtx, exam_dtos.CreateTestCaseInput{
		Name:           "exam_score_case_2_1",
		Input:          []exam_dtos.IOVariableDTO{{Name: "n", Type: string(exam_consts.VariableFormatInt), Value: "2"}},
		ExpectedOutput: exam_dtos.IOVariableDTO{Name: "out", Type: string(exam_consts.VariableFormatInt), Value: "4"},
		IsSample:       false,
		Points:         5,
		ChallengeID:    secondChallengeID,
	})
	if err != nil {
		process.Fail("create second test case 1", err)
	}
	secondTestCaseOneID = secondTestCaseOne.ID

	secondTestCaseTwo, err := process.Application.TestCaseModule.CreateTestCase.Execute(teacherCtx, exam_dtos.CreateTestCaseInput{
		Name:           "exam_score_case_2_2",
		Input:          []exam_dtos.IOVariableDTO{{Name: "n", Type: string(exam_consts.VariableFormatInt), Value: "5"}},
		ExpectedOutput: exam_dtos.IOVariableDTO{Name: "out", Type: string(exam_consts.VariableFormatInt), Value: "10"},
		IsSample:       false,
		Points:         5,
		ChallengeID:    secondChallengeID,
	})
	if err != nil {
		process.Fail("create second test case 2", err)
	}
	secondTestCaseTwoID = secondTestCaseTwo.ID
	process.EndStep()

	process.StartStep("Crear un punto de examen con valor de 70")
	secondExamItem, err := process.Application.ExamItemModule.CreateExamItem.Execute(teacherCtx, exam_dtos.CreateExamItemInput{
		ExamID:      examID,
		ChallengeID: secondChallengeID,
		Order:       2,
		Points:      70,
	})
	if err != nil {
		process.Fail("create second exam item", err)
	}
	secondExamItemID = secondExamItem.ID
	process.EndStep()

	process.StartStep("Iniciar sesión con usuario de estudiante")
	studentAccess := utils.EnsureAuthUserAccess(t, process.Application, studentEmail, password, "Student Test")
	studentCtx = utils.BuildUserCtx(studentAccess)
	studentID = studentAccess.UserData.ID
	process.EndStep()

	process.StartStep("Crear una sesión en el examen")
	createdSession, err := process.Application.SessionModule.CreateSession.Execute(studentCtx, submission_dtos.CreateSessionInput{UserID: studentID, ExamID: examID})
	if err != nil {
		process.Fail("create session", err)
	}
	sessionID = createdSession.ID
	process.EndStep()

	process.StartStep("Crear una revisión")
	createdSubmission, err := process.Application.SubmissionUseCases.CreateSubmission.Execute(studentCtx, submission_dtos.CreateSubmissionInput{
		Code:        "import sys\nprint(int(sys.stdin.read().strip()) * 2)",
		Language:    "python",
		ChallengeID: firstChallengeID,
		SessionID:   sessionID,
	})
	if err != nil {
		process.Fail("create submission", err)
	}
	submissionID = createdSubmission.ID
	process.EndStep()

	process.StartStep("Obtener el status de la revisión hasta que su estado sea accepted o wrong_answer")
	deadline := time.Now().Add(2 * time.Minute)
	var lastStatusOutput *submission_dtos.SubmissionOutputDTO
	timedOut := true
	for time.Now().Before(deadline) {
		statusOutput, statusErr := process.Application.SubmissionUseCases.GetSubmissionStatus.Execute(studentCtx, submission_dtos.GetSubmissionStatusInput{SubmissionID: submissionID})
		if statusErr != nil {
			process.Log(fmt.Sprintf("Polling status error: %v", statusErr))
			time.Sleep(2 * time.Second)
			continue
		}
		lastStatusOutput = statusOutput
		if statusOutput != nil && len(statusOutput.Results) == 3 {
			allTerminal := true
			for _, r := range statusOutput.Results {
				if r.Status != submission_consts.SubmissionStatusAccepted && r.Status != submission_consts.SubmissionStatusWrongAnswer {
					allTerminal = false
					break
				}
			}
			if allTerminal {
				timedOut = false
				break
			}
		}
		time.Sleep(2 * time.Second)
	}
	if timedOut {
		process.Fail("wait scoring terminal status", fmt.Errorf("submission %s did not reach terminal accepted/wrong_answer status before timeout", submissionID))
	}
	process.EndStep()

	process.StartStep("Cerrar la sesión del examen")
	_, err = process.Application.SessionModule.CloseSession.Execute(studentCtx, submission_dtos.CloseSessionInput{SessionID: sessionID})
	if err != nil {
		process.Fail("close session", err)
	}
	process.EndStep()

	process.StartStep("Obtener resultados de su examen como estudiante")
	examScoring, err := process.Application.ScoringModule.GetExamScoring.Execute(studentCtx, submission_dtos.GetExamScoringInput{ExamID: examID, UserID: studentID})
	if err != nil {
		process.Fail("get exam scoring", err)
	}
	if examScoring == nil {
		process.Fail("get exam scoring", fmt.Errorf("expected exam scoring output"))
	}
	process.EndStep()

	process.StartStep("Validar que el resultado del examen sea congruente")
	if lastStatusOutput == nil {
		process.Fail("verify submission score", fmt.Errorf("expected submission status output"))
	}
	if lastStatusOutput.Submission.Score != 50 {
		process.Fail("verify submission score", fmt.Errorf("expected submission score 6, got %d", lastStatusOutput.Submission.Score))
	}
	if examScoring.ExamID != examID {
		process.Fail("verify exam scoring", fmt.Errorf("expected exam id %s, got %s", examID, examScoring.ExamID))
	}
	if examScoring.BestScore != 15 {
		process.Fail("verify exam scoring", fmt.Errorf("expected best score 15, got %d", examScoring.BestScore))
	}
	if len(examScoring.ExamScores) != 1 {
		process.Fail("verify exam scoring", fmt.Errorf("expected 15 exam score, got %d", len(examScoring.ExamScores)))
	}
	if examScoring.ExamScores[0].Score != 15 {
		process.Fail("verify exam scoring", fmt.Errorf("expected stored exam score 15, got %d", examScoring.ExamScores[0].Score))
	}
	process.EndStep()

	process.StartStep("Obtener resultados de los punto del examen como estudiante")
	examScoreDetails, err := process.Application.ScoringModule.GetItemScoring.Execute(studentCtx, submission_dtos.GetItemScoringInput{ExamScoreID: examScoring.ExamScores[0].ID})
	if err != nil {
		process.Fail("get item scoring", err)
	}
	if examScoreDetails == nil {
		process.Fail("get item scoring", fmt.Errorf("expected exam score details output"))
	}
	if examScoreDetails.ExamScoreID != examScoring.ExamScores[0].ID {
		process.Fail("get item scoring", fmt.Errorf("expected exam score id %s, got %s", examScoring.ExamScores[0].ID, examScoreDetails.ExamScoreID))
	}
	if len(examScoreDetails.ExamItems) != 2 {
		process.Fail("get item scoring", fmt.Errorf("expected 2 exam items, got %d", len(examScoreDetails.ExamItems)))
	}
	firstItemFound := false
	secondItemFound := false
	for _, item := range examScoreDetails.ExamItems {
		if item.ExamItem.ID == firstExamItemID {
			firstItemFound = true
			if item.ExamItemScore.Score != 50 {
				process.Fail("get item scoring", fmt.Errorf("expected first exam item score 50, got %d", item.ExamItemScore.Score))
			}
		}
		if item.ExamItem.ID == secondExamItemID {
			secondItemFound = true
			if item.ExamItemScore.Score != 0 {
				process.Fail("get item scoring", fmt.Errorf("expected second exam item score 0, got %d", item.ExamItemScore.Score))
			}
		}
	}
	if !firstItemFound {
		process.Fail("get item scoring", fmt.Errorf("expected first exam item %s", firstExamItemID))
	}
	if !secondItemFound {
		process.Fail("get item scoring", fmt.Errorf("expected second exam item %s", secondExamItemID))
	}
	process.EndStep()

	process.StartStep("Obtener los resultados del examen para cada estudiante como profesor")
	usersScores, err := process.Application.ScoringModule.GetUsersScoresByExam.Execute(teacherCtx, submission_dtos.GetUserScoresByExamInput{ExamID: examID})
	if err != nil {
		process.Fail("get users scores by exam", err)
	}
	if usersScores == nil {
		process.Fail("get users scores by exam", fmt.Errorf("expected users scores output"))
	}
	if usersScores.ExamID != examID {
		process.Fail("get users scores by exam", fmt.Errorf("expected exam id %s, got %s", examID, usersScores.ExamID))
	}
	if len(usersScores.ExamScores) != 1 {
		process.Fail("get users scores by exam", fmt.Errorf("expected 1 user score entry, got %d", len(usersScores.ExamScores)))
	}
	if usersScores.ExamScores[0].UserID != studentID {
		process.Fail("get users scores by exam", fmt.Errorf("expected user id %s, got %s", studentID, usersScores.ExamScores[0].UserID))
	}
	if usersScores.ExamScores[0].ExamScores.BestScore != 15 {
		process.Fail("get users scores by exam", fmt.Errorf("expected best score 15, got %d", usersScores.ExamScores[0].ExamScores.BestScore))
	}
	process.EndStep()

	process.StartStep("Obtener los resultados de todos los exámenes del estudiante como profesor")
	studentExamsScores, err := process.Application.ScoringModule.GetExamsScoresByUser.Execute(teacherCtx, submission_dtos.GetExamScoresByUserInput{UserID: studentID})
	if err != nil {
		process.Fail("get exams scores by user", err)
	}
	if studentExamsScores == nil {
		process.Fail("get exams scores by user", fmt.Errorf("expected exams scores output"))
	}
	if studentExamsScores.UserID != studentID {
		process.Fail("get exams scores by user", fmt.Errorf("expected user id %s, got %s", studentID, studentExamsScores.UserID))
	}
	if len(studentExamsScores.ExamScores) != 1 {
		process.Fail("get exams scores by user", fmt.Errorf("expected 1 exam score entry, got %d", len(studentExamsScores.ExamScores)))
	}
	if studentExamsScores.ExamScores[0].ExamID != examID {
		process.Fail("get exams scores by user", fmt.Errorf("expected exam id %s, got %s", examID, studentExamsScores.ExamScores[0].ExamID))
	}
	if studentExamsScores.PromScore != 15 {
		process.Fail("get exams scores by user", fmt.Errorf("expected prom score 15, got %f", studentExamsScores.PromScore))
	}
	process.EndStep()

	process.End()
}
