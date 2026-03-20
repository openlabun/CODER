package user_repository

import (
	Entities "../../entities/user"
)

type UserRepository interface {
	// SaveUser persists a new user or updates an existing one.
	SaveUser(user *Entities.User) (*Entities.User, error)

	GetUserByID(userID string) (*Entities.User, error)
	GetUserByEmail(email string) (*Entities.User, error)
	GetUserByUsername(username string) (*Entities.User, error)

	ExistsByEmail(email string) (bool, error)
	ExistsByUsername(username string) (bool, error)
}
