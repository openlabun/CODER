package user_factory

import (
	"time"

	constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/user"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	Validations "github.com/openlabun/CODER/apps/api_v2/internal/domain/validations/user"
)

func NewUser (id, username, email, password string) (*Entities.User, error) {
	// Set default role (e.g., student)
	defaultRole := constants.UserRoleStudent

	// Set creation time and update time
	now := time.Now()

	user := &Entities.User{
		ID:        id,
		Username:  username,
		Email:     email,
		Role:      defaultRole, // Default role
		CreatedAt: now,
		LastConnection: now,
	}

	// Validate registration data
	if err := Validations.ValidateRegistrationData(email, username, password); err != nil {
		return nil, err
	}

	return user, nil
}

func ExistingUser (id, username, email string, role constants.UserRole, createdAt, LastConnection time.Time, connected bool) (*Entities.User, error) {
	// Update last connection time if user is connected
	if connected {
		LastConnection = time.Now()
	}

	user := &Entities.User{
		ID:        id,
		Username:  username,
		Email:     email,
		Role:      role,
		CreatedAt: createdAt,
		LastConnection: LastConnection,
	}

	return user, nil
}