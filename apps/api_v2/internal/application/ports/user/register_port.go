package user_ports

import (
	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/user"
)

type RegisterPort interface {
	RegisterUser(email, name, password string) (bool, error)
	RegisterUserDirect(email, password, name string) (*dtos.UserAccess, error)
	VerifyEmail(email, code string) (bool, error)
}