package services

import (
	"backend/dao"
	"backend/domain"
	"errors"
)

type userService struct {
	dao dao.IUserDAO
}

func NewUserService(d dao.IUserDAO) IUserService {
	return &userService{dao: d}
}

func (s *userService) Register(_ domain.RegisterRequest) (*domain.AuthResponse, error) {
	return nil, errors.New("not implemented")
}

func (s *userService) Login(_ domain.LoginRequest) (*domain.AuthResponse, error) {
	return nil, errors.New("not implemented")
}
