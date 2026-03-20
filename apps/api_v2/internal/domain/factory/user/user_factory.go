package user_factory

import (
	"time"
	Entities "../../entities/user"
	Validations "../../validations/user"
)

func NewUser (id, username, email, password string) (*Entities.User, error) {
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
		UpdatedAt: now,
	}

	// Validate registration data
	if err := Validations.ValidateRegistrationData(email, username, password); err != nil {
		return nil, err
	}

	return user, nil
}

func ExistingUser (id, username, email, password string, role Entities.UserRole, createdAt, updatedAt time.Time) (*Entities.User, error) {
	// Set update time to now
	now := time.Now()
	
	user := &Entities.User{
		ID:        id,
		Username:  username,
		Email:     email,
		Role:      role,
		CreatedAt: createdAt,
		UpdatedAt: now,
	}

	return user, nil
}