package dtos

import (
	Entities "../../../domain/entities/user"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserAccess struct {
	Token 	 *Token `json:"token,omitempty"`
	UserData *Entities.User `json:"user_data,omitempty"`
}