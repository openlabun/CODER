package user_ports

import (
	Entities "../../../domain/entities/user"
	dtos "../../dtos/user"
)

// UserServicePort abstracts user persistence/query operations needed by user use cases.
type UserServicePort interface {
	GetUserByID(userID string) (*Entities.User, error)
	GetUserByEmail(email string) (*Entities.User, error)
	GetUserByUsername(username string) (*Entities.User, error)
	ExistsByEmail(email string) (bool, error)
	ExistsByUsername(username string) (bool, error)
}

// PasswordHasherPort abstracts password hashing/verification.
type PasswordHasherPort interface {
	Hash(plain string) (string, error)
	Compare(hash, plain string) (bool, error)
}

// TokenServicePort abstracts token generation/refresh logic.
type TokenServicePort interface {
	RefreshUserToken(refreshToken string) (*dtos.Token, error)
}
