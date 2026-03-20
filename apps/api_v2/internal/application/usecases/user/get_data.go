package user_usecases

import (
	"fmt"

	Entities "../../../domain/entities/user"
	ports "../../ports/user"
)

type GetDataUseCase struct {
	userService ports.LoginPort
}

func NewGetDataUseCase(userService ports.LoginPort) *GetDataUseCase {
	return &GetDataUseCase{userService: userService}
}

func (uc *GetDataUseCase) Execute(userID string) (*Entities.User, error) {

	// Validate if user exists by ID
	user, err := uc.userService.GetUserData(userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching user: %w", err)
	}

	return user, nil
}
