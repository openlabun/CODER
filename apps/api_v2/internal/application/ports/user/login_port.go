package user_ports

import (
	dtos "../../dtos/user"
	Entities "../../../domain/entities/user"
)

type LoginPort interface {
	LoginUser(email, password string) (*dtos.UserAccess, error)
	GetUserData(userID string) (*Entities.User, error)
}