package user_usecases

import (
	"context"
	"fmt"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	ports "github.com/openlabun/CODER/apps/api_v2/internal/application/ports/user"
)

type GetDataUseCase struct {
	userService ports.LoginPort
}

func NewGetDataUseCase(userService ports.LoginPort) *GetDataUseCase {
	return &GetDataUseCase{userService: userService}
}

func (uc *GetDataUseCase) Execute(ctx context.Context, userID string) (*Entities.User, error) {

	// Validate if user exists by ID
	user, err := uc.userService.GetUserData(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching user: %w", err)
	}

	return user, nil
}
