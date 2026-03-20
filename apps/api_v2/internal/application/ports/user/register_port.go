package user_ports

import (
	dtos "../../dtos/user"
)

type RegisterPort interface {
	RegisterUser(email, name, password string) (*dtos.UserAccess, error)
}