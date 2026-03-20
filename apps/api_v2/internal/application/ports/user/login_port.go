package user_ports

import (
	dtos "../../dtos/user"
	Entities "../../../domain/entities/user"
)

type LoginPort interface {
	LoginUser(email, password string) (*dtos.UserAccess, error)
	GetUserData(email string) (*Entities.User, error)
}