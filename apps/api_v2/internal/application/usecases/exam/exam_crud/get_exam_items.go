package exam_usecases

import (
	"context"
	"fmt"

	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	
	user_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	repositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"

	mapper "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam/mapper"
	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
)

type GetExamItemsUseCase struct {
	challengeRepository repositories.ChallengeRepository
	userRepository      userRepository.UserRepository
	examRepository      examRepository.ExamRepository
	examItemRepository  examRepository.ExamItemRepository
}

func NewGetExamItemsUseCase(challengeRepository repositories.ChallengeRepository, userRepository userRepository.UserRepository, examRepository examRepository.ExamRepository, examItemRepository examRepository.ExamItemRepository) *GetExamItemsUseCase {
	return &GetExamItemsUseCase{challengeRepository: challengeRepository, userRepository: userRepository, examRepository: examRepository, examItemRepository: examItemRepository}
}

func (uc *GetExamItemsUseCase) Execute(ctx context.Context, input dtos.GetExamItemsInput) ([]dtos.ExamItemDTO, error) {
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

	// [STEP 2] Verify exam exists
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

	// [STEP 4] Get challenges of the exam
	challenges, err := uc.challengeRepository.GetChallengesByExamID(ctx, input.ExamID)
	if err != nil {
		return nil, err
	}	

	// [STEP 5] Get exam items details and return them
	examItems, err := uc.examItemRepository.GetExamItem(ctx, &input.ExamID, nil)
	if err != nil {
		return nil, err
	}

	// [STEP 6] Map exam items and challenges to exam item dtos
	examItemDTOs, err := MapExamItemDTOs(examItems, challenges)
	if err != nil {
		return nil, err
	}

	// [STEP 5] Return challenge details
	return examItemDTOs, nil
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

func MapExamItemDTOs(examItems []*Entities.ExamItem, challenges []*Entities.Challenge) ([]dtos.ExamItemDTO, error) {
	var examItemDTOs []dtos.ExamItemDTO
	for _, examItem := range examItems {
		var challenge *Entities.Challenge
		for _, c := range challenges {
			if c.ID == examItem.ChallengeID {
				challenge = c
				break
			}
		}
		if challenge == nil {
			return nil, fmt.Errorf("challenge not found for exam item with id %q", examItem.ID)
		}
		dto, err := mapper.MapExamItemDTO(examItem, challenge)
		if err != nil {
			return nil, err
		}

		if dto == nil {
			return nil, fmt.Errorf("error mapping exam item with id %q", examItem.ID)
		}
		examItemDTOs = append(examItemDTOs, *dto)
	}
	return examItemDTOs, nil
}