package services

import (
	"backend/dao"
	"backend/domain"
	"backend/utils"
	"errors"

	"gorm.io/gorm"
)

type userService struct {
	dao dao.IUserDAO
}

func NewUserService(d dao.IUserDAO) IUserService {
	return &userService{dao: d}
}

func (s *userService) Register(req domain.RegisterRequest) (*domain.AuthResponse, error) {
	_, err := s.dao.FindByEmail(req.Email)
	if err == nil {
		return nil, errors.New("el email ya está registrado")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	salt := utils.GenerateSalt()
	user := &domain.User{
		Nombre:       req.Nombre,
		Email:        req.Email,
		PasswordHash: utils.HashPassword(req.Password, salt),
		PasswordSalt: salt,
		Rol:          "cliente",
	}
	if err := s.dao.Create(user); err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(user.ID, user.Rol)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		Token: token,
		User: domain.UserResponse{
			ID:     user.ID,
			Nombre: user.Nombre,
			Email:  user.Email,
			Rol:    user.Rol,
		},
	}, nil
}

func (s *userService) Login(req domain.LoginRequest) (*domain.AuthResponse, error) {
	user, err := s.dao.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	if utils.HashPassword(req.Password, user.PasswordSalt) != user.PasswordHash {
		return nil, errors.New("credenciales inválidas")
	}

	token, err := utils.GenerateToken(user.ID, user.Rol)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		Token: token,
		User: domain.UserResponse{
			ID:     user.ID,
			Nombre: user.Nombre,
			Email:  user.Email,
			Rol:    user.Rol,
		},
	}, nil
}
