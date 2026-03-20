package user_repository

import (
	Entities "../../entities/user"
)

type UserRepository interface {
	SaveUser(user *Entities.User) (*Entities.User, error)

	ValidateCredentials(email, password string) (*Entities.User, error)
	refreshToken(refreshToken string) (string, error)
}

