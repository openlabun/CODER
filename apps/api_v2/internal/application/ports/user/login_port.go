package user_ports

import (
	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/user"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
)

type LoginPort interface {
	LoginUser(email, password string) (*dtos.UserAccess, error)
	GetUserData(email string) (*Entities.User, error)
}