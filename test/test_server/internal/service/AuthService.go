package service

import (
	"errors"
	"serverEngineTests/internal/model"
)

type AuthService struct {
}

// @service
func NewAuthService() *AuthService {
	return &AuthService{}
}

func (a *AuthService) Login(username, password string) (*model.User, error) {
	if username == password {
		return &model.User{
			Username: username,
			Password: password,
			Roles:    []string{"admin"},
		}, nil
	}
	return nil, errors.New("error")
}
