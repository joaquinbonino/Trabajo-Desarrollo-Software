package controllers_test

import (
	"backend/controllers"
	"backend/domain"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockWaitlistService struct{ mock.Mock }

func (m *mockWaitlistService) JoinWaitlist(userID, eventID uint) error {
	args := m.Called(userID, eventID)
	return args.Error(0)
}
func (m *mockWaitlistService) GetMyWaitlist(userID uint) ([]domain.WaitlistResponse, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.WaitlistResponse), args.Error(1)
}

func setupWaitlistRouter(ctrl *controllers.WaitlistController, userID uint) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(injectUser(userID))
	r.POST("/events/:id/waitlist", ctrl.Join)
	r.GET("/waitlist/mine", ctrl.GetMine)
	return r
}

func TestJoinWaitlistEndpoint_Exitoso(t *testing.T) {
	svc := new(mockWaitlistService)
	svc.On("JoinWaitlist", uint(42), uint(1)).Return(nil)

	r := setupWaitlistRouter(controllers.NewWaitlistController(svc), 42)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/events/1/waitlist", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestJoinWaitlistEndpoint_HayCupo(t *testing.T) {
	svc := new(mockWaitlistService)
	svc.On("JoinWaitlist", uint(42), uint(1)).
		Return(errors.New("todavía hay entradas disponibles, comprá directamente"))

	r := setupWaitlistRouter(controllers.NewWaitlistController(svc), 42)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/events/1/waitlist", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJoinWaitlistEndpoint_IDInvalido(t *testing.T) {
	svc := new(mockWaitlistService)

	r := setupWaitlistRouter(controllers.NewWaitlistController(svc), 42)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/events/abc/waitlist", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetMyWaitlistEndpoint_Exitoso(t *testing.T) {
	svc := new(mockWaitlistService)
	svc.On("GetMyWaitlist", uint(42)).Return([]domain.WaitlistResponse{
		{ID: 1, Estado: "pendiente", Event: domain.EventResponse{ID: 1, Titulo: "Recital"}},
	}, nil)

	r := setupWaitlistRouter(controllers.NewWaitlistController(svc), 42)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/waitlist/mine", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotNil(t, resp["data"])
}

func TestGetMyWaitlistEndpoint_Error(t *testing.T) {
	svc := new(mockWaitlistService)
	svc.On("GetMyWaitlist", uint(42)).Return(nil, errors.New("db error"))

	r := setupWaitlistRouter(controllers.NewWaitlistController(svc), 42)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/waitlist/mine", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
