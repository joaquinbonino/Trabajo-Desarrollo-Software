package controllers_test

import (
	"backend/controllers"
	"backend/domain"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockEventService struct{ mock.Mock }

func (m *mockEventService) ListEvents(categoria string) ([]domain.EventResponse, error) {
	args := m.Called(categoria)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.EventResponse), args.Error(1)
}
func (m *mockEventService) GetEvent(id uint) (*domain.EventResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.EventResponse), args.Error(1)
}
func (m *mockEventService) CreateEvent(req domain.CreateEventRequest) (*domain.EventResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.EventResponse), args.Error(1)
}
func (m *mockEventService) UpdateEvent(id uint, req domain.UpdateEventRequest) (*domain.EventResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.EventResponse), args.Error(1)
}
func (m *mockEventService) CancelEvent(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func setupEventRouter(ctrl *controllers.EventController) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/events", ctrl.List)
	r.GET("/events/:id", ctrl.Get)
	return r
}

func sampleEventResponse() domain.EventResponse {
	return domain.EventResponse{
		ID:             1,
		Titulo:         "Concierto",
		Categoria:      "Música",
		FechaHora:      time.Now().Add(24 * time.Hour),
		CapacidadTotal: 100,
	}
}

func TestListEventsEndpoint_Exitoso(t *testing.T) {
	svc := new(mockEventService)
	svc.On("ListEvents", "").Return([]domain.EventResponse{sampleEventResponse()}, nil)

	r := setupEventRouter(controllers.NewEventController(svc))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/events", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotNil(t, resp["data"])
}

func TestListEventsEndpoint_ConFiltro(t *testing.T) {
	svc := new(mockEventService)
	svc.On("ListEvents", "Música").Return([]domain.EventResponse{sampleEventResponse()}, nil)

	r := setupEventRouter(controllers.NewEventController(svc))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/events?categoria=Música", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListEventsEndpoint_ErrorServicio(t *testing.T) {
	svc := new(mockEventService)
	svc.On("ListEvents", "").Return(nil, errors.New("db error"))

	r := setupEventRouter(controllers.NewEventController(svc))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/events", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetEventEndpoint_Exitoso(t *testing.T) {
	ev := sampleEventResponse()
	svc := new(mockEventService)
	svc.On("GetEvent", uint(1)).Return(&ev, nil)

	r := setupEventRouter(controllers.NewEventController(svc))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/events/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotNil(t, resp["data"])
}

func TestGetEventEndpoint_NoExiste(t *testing.T) {
	svc := new(mockEventService)
	svc.On("GetEvent", uint(99)).Return(nil, errors.New("evento no encontrado"))

	r := setupEventRouter(controllers.NewEventController(svc))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/events/99", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetEventEndpoint_IDInvalido(t *testing.T) {
	svc := new(mockEventService)

	r := setupEventRouter(controllers.NewEventController(svc))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/events/abc", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
