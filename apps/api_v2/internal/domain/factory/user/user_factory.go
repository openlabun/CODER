package user_factory

import (
	"time"
	Entities "../../entities/user"
	Validations "../../validations/user"
)

func NewUser (id, username, email string) (*Entities.User, error) {
	// Set default role (e.g., student)
	defaultRole := Entities.UserRoleStudent

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

func ExistingUser (id, username, email string, role Entities.UserRole, createdAt, LastConnection time.Time, connected bool) (*Entities.User, error) {
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