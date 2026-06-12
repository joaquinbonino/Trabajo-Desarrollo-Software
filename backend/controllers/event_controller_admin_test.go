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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// setupAdminEventRouter registra los endpoints de administración de eventos.
// Usa /admin/events para ListAll para no chocar con el GET /events público.
func setupAdminEventRouter(ctrl *controllers.EventController) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/events", ctrl.Create)
	r.PUT("/events/:id", ctrl.Update)
	r.GET("/admin/events", ctrl.ListAll)
	r.GET("/events/:id/report", ctrl.Report)
	r.PATCH("/events/:id/cancel", ctrl.Cancel)
	return r
}

// Create

func TestCreateEventEndpoint_Exitoso(t *testing.T) {
	ev := sampleEventResponse()
	svc := new(mockEventService)
	svc.On("CreateEvent", mock.AnythingOfType("domain.CreateEventRequest")).Return(&ev, nil)

	r := setupAdminEventRouter(controllers.NewEventController(svc))
	body, _ := json.Marshal(map[string]interface{}{
		"titulo":          "Concierto",
		"fecha_hora":      time.Now().Add(24 * time.Hour),
		"capacidad_total": 100,
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/events", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateEventEndpoint_BodyInvalido(t *testing.T) {
	svc := new(mockEventService)

	r := setupAdminEventRouter(controllers.NewEventController(svc))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/events", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateEventEndpoint_ErrorServicio(t *testing.T) {
	svc := new(mockEventService)
	svc.On("CreateEvent", mock.AnythingOfType("domain.CreateEventRequest")).
		Return(nil, errors.New("db error"))

	r := setupAdminEventRouter(controllers.NewEventController(svc))
	body, _ := json.Marshal(map[string]interface{}{
		"titulo":          "Concierto",
		"fecha_hora":      time.Now().Add(24 * time.Hour),
		"capacidad_total": 100,
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/events", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// Update

func TestUpdateEventEndpoint_Exitoso(t *testing.T) {
	ev := sampleEventResponse()
	svc := new(mockEventService)
	svc.On("UpdateEvent", uint(1), mock.AnythingOfType("domain.UpdateEventRequest")).
		Return(&ev, nil)

	r := setupAdminEventRouter(controllers.NewEventController(svc))
	body, _ := json.Marshal(map[string]interface{}{"titulo": "Nuevo"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/events/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateEventEndpoint_IDInvalido(t *testing.T) {
	svc := new(mockEventService)

	r := setupAdminEventRouter(controllers.NewEventController(svc))
	body, _ := json.Marshal(map[string]interface{}{"titulo": "Nuevo"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/events/abc", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateEventEndpoint_BodyInvalido(t *testing.T) {
	svc := new(mockEventService)

	r := setupAdminEventRouter(controllers.NewEventController(svc))
	w := httptest.NewRecorder()
	// capacidad_total negativa viola el binding omitempty,min=1
	req, _ := http.NewRequest(http.MethodPut, "/events/1", bytes.NewBufferString(`{"capacidad_total": -5}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateEventEndpoint_ErrorNegocio(t *testing.T) {
	svc := new(mockEventService)
	svc.On("UpdateEvent", uint(1), mock.AnythingOfType("domain.UpdateEventRequest")).
		Return(nil, errors.New("no se puede modificar un evento cancelado"))

	r := setupAdminEventRouter(controllers.NewEventController(svc))
	body, _ := json.Marshal(map[string]interface{}{"titulo": "Nuevo"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/events/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ListAll

func TestListAllEventsEndpoint_Exitoso(t *testing.T) {
	svc := new(mockEventService)
	svc.On("ListAllEvents").Return([]domain.EventResponse{sampleEventResponse()}, nil)

	r := setupAdminEventRouter(controllers.NewEventController(svc))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/admin/events", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotNil(t, resp["data"])
}

func TestListAllEventsEndpoint_ErrorServicio(t *testing.T) {
	svc := new(mockEventService)
	svc.On("ListAllEvents").Return(nil, errors.New("db error"))

	r := setupAdminEventRouter(controllers.NewEventController(svc))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/admin/events", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// Report

func TestReportEndpoint_Exitoso(t *testing.T) {
	report := &domain.EventReportResponse{
		EventID:             1,
		Titulo:              "Concierto",
		CapacidadTotal:      100,
		EntradasVendidas:    10,
		EntradasDisponibles: 90,
	}
	svc := new(mockEventService)
	svc.On("GetEventReport", uint(1)).Return(report, nil)

	r := setupAdminEventRouter(controllers.NewEventController(svc))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/events/1/report", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotNil(t, resp["data"])
}

func TestReportEndpoint_NoExiste(t *testing.T) {
	svc := new(mockEventService)
	svc.On("GetEventReport", uint(99)).Return(nil, errors.New("evento no encontrado"))

	r := setupAdminEventRouter(controllers.NewEventController(svc))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/events/99/report", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestReportEndpoint_IDInvalido(t *testing.T) {
	svc := new(mockEventService)

	r := setupAdminEventRouter(controllers.NewEventController(svc))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/events/abc/report", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Cancel

func TestCancelEventEndpoint_Exitoso(t *testing.T) {
	svc := new(mockEventService)
	svc.On("CancelEvent", uint(1)).Return(nil)

	r := setupAdminEventRouter(controllers.NewEventController(svc))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/events/1/cancel", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCancelEventEndpoint_IDInvalido(t *testing.T) {
	svc := new(mockEventService)

	r := setupAdminEventRouter(controllers.NewEventController(svc))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/events/abc/cancel", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCancelEventEndpoint_ErrorNegocio(t *testing.T) {
	svc := new(mockEventService)
	svc.On("CancelEvent", uint(1)).Return(errors.New("el evento ya está cancelado"))

	r := setupAdminEventRouter(controllers.NewEventController(svc))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/events/1/cancel", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
