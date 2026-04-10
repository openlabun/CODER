package exam_usecases

import (
	"context"
	"fmt"
	"time"

	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	
	user_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
)

type GetExamScores struct {
	challengeRepository examRepository.ChallengeRepository
	examRepository 		examRepository.ExamRepository
	userRepository      userRepository.UserRepository
}

func NewGetExamScores(challengeRepository examRepository.ChallengeRepository, examRepository examRepository.ExamRepository, userRepository userRepository.UserRepository) *GetExamScores {
	return &GetExamScores{challengeRepository: challengeRepository, examRepository: examRepository, userRepository: userRepository}
}

func (uc *GetExamScores) Execute(ctx context.Context, input dtos.GetExamScoreInput) ([]*dtos.ExamScoreDTO, error) {
	// [STEP 1] Verify user and get its role
	userEmail, err := services.UserEmailFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepository.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// [STEP 2] Verify if user has permissions to view the exam scores
	if input.UserID != nil && user.ID != *input.UserID && user.Role == user_entities.UserRoleStudent {
		return nil, fmt.Errorf("students are not allowed to access other users' exam scores")
	}

	// [STEP 3] Get Exam and validate it exists
	exam, err := uc.examRepository.GetExamByID(ctx, input.ExamID)
	if err != nil {
		return nil, err
	}

	if exam == nil {
		return nil, fmt.Errorf("exam with id %q not found", input.ExamID)
	}

	// [STEP 4] Validate if user is teacher, only can get scores if is the owner of the exam or exam is "public" or "teachers"
	if user.Role == user_entities.UserRoleProfessor {
		if exam.ProfessorID != user.ID && exam.Visibility != Entities.VisibilityPublic && exam.Visibility != Entities.VisibilityTeachers {
			return nil, fmt.Errorf("professors can only access scores of their own exams or public exams")
		}
	}

	// [STEP 5] Get scores for the exam (Mockup)
	test_user_id := "c585dd74-b0d5-4638-9657-bccf7ab37acd"
	// Mockup exam
	start := time.Now().Add(-1 * time.Hour)
	end := time.Now().Add(1 * time.Hour)
	exam = &Entities.Exam{
		ID:          input.ExamID,
		Title:       "Sample Exam",
		Description: "This is a sample exam",
		Visibility:  Entities.VisibilityPublic,
		StartTime:   start,
		EndTime:     &end,
		AllowLateSubmissions: false,
		TimeLimit:   60,
		TryLimit:    1,
		ProfessorID: "professor123",
	}

	// IOVariable
	inputVar := &Entities.IOVariable{
		Name: "input1",
		Type: "string",
		Value: "hola",
	}
	inputVars := []Entities.IOVariable{*inputVar}

	outputVar := &Entities.IOVariable{
		Name: "output1",
		Type: "string",
		Value: "hola",
	}

	// Mockup challenges
	challenge_1 := &Entities.Challenge{
		ID: "challenge1",
		Title: "Challenge 1",
		Description: "This is the first challenge",
		Status: Entities.ChallengeStatusPublished,
		Tags: []string{"tag1", "tag2"},
		Difficulty: Entities.ChallengeDifficultyHard,
		WorkerTimeLimit: 10,
		WorkerMemoryLimit: 10,
		InputVariables: inputVars,
		OutputVariable: *outputVar,
		Constraints: "tst",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: test_user_id,
	}

	challenge_2 := &Entities.Challenge{
		ID: "challenge2",
		Title: "Challenge 2",
		Description: "This is the second challenge",
		Status: Entities.ChallengeStatusPublished,
		Tags: []string{"tag1", "tag2"},
		Difficulty: Entities.ChallengeDifficultyHard,
		WorkerTimeLimit: 10,
		WorkerMemoryLimit: 10,
		InputVariables: inputVars,
		OutputVariable: *outputVar,
		Constraints: "tst",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: test_user_id,
	}

	// Mockup exam item scores
	examItemScores := []*dtos.ExamItemScoreDTO{
		{
			ExamItemID: "item1",
			Challenge:  challenge_1,
			Score:      85.0,
		},
		{
			ExamItemID: "item2",
			Challenge:  challenge_2,
			Score:      90.0,
		},
	}

	// Mockup exam score
	examScore := &dtos.ExamScoreDTO{
		UserID: test_user_id,
		Exam: exam,
		Score: 87.5,
		ExamItemScores: examItemScores,
	}

	examScores := []*dtos.ExamScoreDTO{examScore}

	return examScores, nil
}