package services_test

import (
	"backend/domain"
	"backend/services"
	"backend/utils"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type mockUserDAO struct{ mock.Mock }

func (m *mockUserDAO) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}
func (m *mockUserDAO) FindByID(id uint) (*domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.User), args.Error(1)
}
func (m *mockUserDAO) FindByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}
func (m *mockUserDAO) Update(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestRegister_Exitoso(t *testing.T) {
	d := new(mockUserDAO)
	d.On("FindByEmail", "test@mail.com").Return(nil, gorm.ErrRecordNotFound)
	d.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)

	svc := services.NewUserService(d)
	resp, err := svc.Register(domain.RegisterRequest{
		Nombre:   "Juan",
		Email:    "test@mail.com",
		Password: "password123",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, "test@mail.com", resp.User.Email)
	assert.Equal(t, "cliente", resp.User.Rol)
}

func TestRegister_EmailDuplicado(t *testing.T) {
	d := new(mockUserDAO)
	d.On("FindByEmail", "test@mail.com").Return(&domain.User{Email: "test@mail.com"}, nil)

	svc := services.NewUserService(d)
	resp, err := svc.Register(domain.RegisterRequest{
		Nombre:   "Juan",
		Email:    "test@mail.com",
		Password: "password123",
	})

	assert.Nil(t, resp)
	assert.EqualError(t, err, "el email ya está registrado")
}

func TestLogin_Exitoso(t *testing.T) {
	salt := "testsalt"
	d := new(mockUserDAO)
	d.On("FindByEmail", "test@mail.com").Return(&domain.User{
		Email:        "test@mail.com",
		PasswordHash: utils.HashPassword("password123", salt),
		PasswordSalt: salt,
		Rol:          "cliente",
	}, nil)

	svc := services.NewUserService(d)
	resp, err := svc.Login(domain.LoginRequest{
		Email:    "test@mail.com",
		Password: "password123",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, "test@mail.com", resp.User.Email)
}

func TestLogin_PasswordIncorrecta(t *testing.T) {
	d := new(mockUserDAO)
	d.On("FindByEmail", "test@mail.com").Return(&domain.User{
		Email:        "test@mail.com",
		PasswordHash: "hashincorrecto",
		PasswordSalt: "salt",
		Rol:          "cliente",
	}, nil)

	svc := services.NewUserService(d)
	resp, err := svc.Login(domain.LoginRequest{
		Email:    "test@mail.com",
		Password: "wrongpassword",
	})

	assert.Nil(t, resp)
	assert.EqualError(t, err, "credenciales inválidas")
}

func TestLogin_UsuarioNoExiste(t *testing.T) {
	d := new(mockUserDAO)
	d.On("FindByEmail", "noexiste@mail.com").Return(nil, errors.New("not found"))

	svc := services.NewUserService(d)
	resp, err := svc.Login(domain.LoginRequest{
		Email:    "noexiste@mail.com",
		Password: "password123",
	})

	assert.Nil(t, resp)
	assert.EqualError(t, err, "credenciales inválidas")
}

func TestRegister_ErrorAlCrear(t *testing.T) {
	d := new(mockUserDAO)
	d.On("FindByEmail", "test@mail.com").Return(nil, gorm.ErrRecordNotFound)
	d.On("Create", mock.AnythingOfType("*domain.User")).Return(errors.New("db error"))

	svc := services.NewUserService(d)
	resp, err := svc.Register(domain.RegisterRequest{
		Nombre:   "Juan",
		Email:    "test@mail.com",
		Password: "password123",
	})

	assert.Nil(t, resp)
	assert.Error(t, err)
}

func TestRegister_ErrorInesperadoEnBusqueda(t *testing.T) {
	d := new(mockUserDAO)
	// Un error que NO es "registro no encontrado" debe propagarse tal cual.
	d.On("FindByEmail", "test@mail.com").Return(nil, errors.New("db caída"))

	svc := services.NewUserService(d)
	resp, err := svc.Register(domain.RegisterRequest{
		Nombre:   "Juan",
		Email:    "test@mail.com",
		Password: "password123",
	})

	assert.Nil(t, resp)
	assert.EqualError(t, err, "db caída")
}
