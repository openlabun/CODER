package user_ports

import (
	dtos "../../dtos/user"
)

type RegisterPort interface {
	RegisterUser(email, name, password string) (bool, error)
	RegisterUserDirect(email, password, name string) (*dtos.UserAccess, error)
	VerifyEmail(email, code string) (bool, error)
}