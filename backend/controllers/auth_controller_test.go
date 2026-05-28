package controllers_test

import (
	"backend/controllers"
	"backend/domain"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserService struct{ mock.Mock }

func (m *mockUserService) Register(req domain.RegisterRequest) (*domain.AuthResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AuthResponse), args.Error(1)
}
func (m *mockUserService) Login(req domain.LoginRequest) (*domain.AuthResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AuthResponse), args.Error(1)
}

func setupRouter(ctrl *controllers.AuthController) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/auth/register", ctrl.Register)
	r.POST("/auth/login", ctrl.Login)
	return r
}

func TestRegisterEndpoint_Exitoso(t *testing.T) {
	svc := new(mockUserService)
	svc.On("Register", mock.AnythingOfType("domain.RegisterRequest")).Return(
		&domain.AuthResponse{Token: "jwt-token", User: domain.UserResponse{Email: "test@mail.com"}}, nil,
	)

	r := setupRouter(controllers.NewAuthController(svc))
	body, _ := json.Marshal(map[string]string{
		"nombre": "Juan", "email": "test@mail.com", "password": "password123",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotNil(t, resp["data"])
}

func TestRegisterEndpoint_EmailDuplicado(t *testing.T) {
	svc := new(mockUserService)
	svc.On("Register", mock.AnythingOfType("domain.RegisterRequest")).Return(
		nil, errors.New("el email ya está registrado"),
	)

	r := setupRouter(controllers.NewAuthController(svc))
	body, _ := json.Marshal(map[string]string{
		"nombre": "Juan", "email": "test@mail.com", "password": "password123",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestRegisterEndpoint_BodyInvalido(t *testing.T) {
	svc := new(mockUserService)
	r := setupRouter(controllers.NewAuthController(svc))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLoginEndpoint_Exitoso(t *testing.T) {
	svc := new(mockUserService)
	svc.On("Login", mock.AnythingOfType("domain.LoginRequest")).Return(
		&domain.AuthResponse{Token: "jwt-token", User: domain.UserResponse{Email: "test@mail.com"}}, nil,
	)

	r := setupRouter(controllers.NewAuthController(svc))
	body, _ := json.Marshal(map[string]string{
		"email": "test@mail.com", "password": "password123",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLoginEndpoint_CredencialesInvalidas(t *testing.T) {
	svc := new(mockUserService)
	svc.On("Login", mock.AnythingOfType("domain.LoginRequest")).Return(
		nil, errors.New("credenciales inválidas"),
	)

	r := setupRouter(controllers.NewAuthController(svc))
	body, _ := json.Marshal(map[string]string{
		"email": "test@mail.com", "password": "wrongpassword",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
