package exam_usecases

import (
	"context"
	"fmt"

	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	
	user_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	repositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
)

type GetChallengesByExamUseCase struct {
	challengeRepository repositories.ChallengeRepository
	userRepository      userRepository.UserRepository
	examRepository      examRepository.ExamRepository
}

func NewGetChallengesByExamUseCase(challengeRepository repositories.ChallengeRepository, userRepository userRepository.UserRepository, examRepository examRepository.ExamRepository) *GetChallengesByExamUseCase {
	return &GetChallengesByExamUseCase{challengeRepository: challengeRepository, userRepository: userRepository, examRepository: examRepository}
}

func (uc *GetChallengesByExamUseCase) Execute(ctx context.Context, input dtos.GetChallengesByExamInput) ([]*Entities.Challenge, error) {
	// [STEP 1] Verify user and get its role
	userEmail, err := services.UserEmailFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepository.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return nil, err
	}

	role := user.Role

	// [STEP 2] Verify challenge exists
	challenges, err := uc.challengeRepository.GetChallengesByExamID(ctx, input.ExamID)
	if err != nil {
		return nil, err
	}

	// [STEP 3] Verify exam exists
	exam, err := uc.examRepository.GetExamByID(ctx, input.ExamID)
	if err != nil {
		return nil, err
	}
	
	if exam == nil {
		return nil, fmt.Errorf("exam with id %q does not exist", input.ExamID)
	}
	

	// [STEP 3] If user is teacher Verify that exam belongs to the teacher or exam is public/teachers
	if role == user_entities.UserRoleProfessor {
		if exam.ProfessorID != user.ID && exam.Visibility != Entities.VisibilityPublic && exam.Visibility != Entities.VisibilityTeachers {
			return nil, fmt.Errorf("user does not have permissions to access this exam")
		}
	}

	// [STEP 4] If user is student verify that exam is accessible and filter challenges to only return published ones
	if role == user_entities.UserRoleStudent {
		if exam.Visibility == Entities.VisibilityPrivate || exam.Visibility == Entities.VisibilityTeachers {
			return nil, fmt.Errorf("user does not have permissions to access this exam")
		}

		challenges = filterPublishedChallenges(challenges)
	}

	// [STEP 5] Return challenge details
	return challenges, nil
}

func filterPublishedChallenges(challenges []*Entities.Challenge) []*Entities.Challenge {
	visibleChallenges := []*Entities.Challenge{}
	for _, challenge := range challenges {
		if challenge.Status == Entities.ChallengeStatusPublished {
			visibleChallenges = append(visibleChallenges, challenge)
		}
	}

	return visibleChallenges
}