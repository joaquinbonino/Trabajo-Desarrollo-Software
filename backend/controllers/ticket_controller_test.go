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

type mockTicketService struct{ mock.Mock }

func (m *mockTicketService) BuyTicket(userID uint, req domain.BuyTicketRequest) (*domain.TicketResponse, error) {
	args := m.Called(userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TicketResponse), args.Error(1)
}
func (m *mockTicketService) GetMyTickets(userID uint) ([]domain.TicketResponse, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.TicketResponse), args.Error(1)
}
func (m *mockTicketService) CancelTicket(ticketID, userID uint) error {
	args := m.Called(ticketID, userID)
	return args.Error(0)
}
func (m *mockTicketService) TransferTicket(ticketID, userID uint, req domain.TransferTicketRequest) error {
	args := m.Called(ticketID, userID, req)
	return args.Error(0)
}

// injectUser simula el middleware de auth inyectando el user_id en el contexto
func injectUser(userID uint) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	}
}

func setupTicketRouter(ctrl *controllers.TicketController, userID uint) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(injectUser(userID))
	r.POST("/tickets", ctrl.Buy)
	r.GET("/tickets/mine", ctrl.GetMine)
	r.DELETE("/tickets/:id", ctrl.Cancel)
	r.PUT("/tickets/:id/transfer", ctrl.Transfer)
	return r
}

func sampleTicketResponse() *domain.TicketResponse {
	return &domain.TicketResponse{
		ID:          1,
		Estado:      "activo",
		FechaCompra: time.Now(),
		Event:       domain.EventResponse{ID: 1, Titulo: "Concierto"},
	}
}

func TestBuyTicketEndpoint_Exitoso(t *testing.T) {
	svc := new(mockTicketService)
	svc.On("BuyTicket", uint(42), domain.BuyTicketRequest{EventID: 1}).
		Return(sampleTicketResponse(), nil)

	r := setupTicketRouter(controllers.NewTicketController(svc), 42)
	body, _ := json.Marshal(map[string]int{"event_id": 1})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/tickets", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestBuyTicketEndpoint_SinCupo(t *testing.T) {
	svc := new(mockTicketService)
	svc.On("BuyTicket", uint(42), domain.BuyTicketRequest{EventID: 1}).
		Return(nil, errors.New("no hay entradas disponibles"))

	r := setupTicketRouter(controllers.NewTicketController(svc), 42)
	body, _ := json.Marshal(map[string]int{"event_id": 1})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/tickets", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBuyTicketEndpoint_BodyInvalido(t *testing.T) {
	svc := new(mockTicketService)

	r := setupTicketRouter(controllers.NewTicketController(svc), 42)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/tickets", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetMineEndpoint_Exitoso(t *testing.T) {
	svc := new(mockTicketService)
	svc.On("GetMyTickets", uint(42)).Return([]domain.TicketResponse{*sampleTicketResponse()}, nil)

	r := setupTicketRouter(controllers.NewTicketController(svc), 42)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/tickets/mine", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotNil(t, resp["data"])
}

func TestCancelTicketEndpoint_Exitoso(t *testing.T) {
	svc := new(mockTicketService)
	svc.On("CancelTicket", uint(1), uint(42)).Return(nil)

	r := setupTicketRouter(controllers.NewTicketController(svc), 42)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/tickets/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCancelTicketEndpoint_NoPermiso(t *testing.T) {
	svc := new(mockTicketService)
	svc.On("CancelTicket", uint(1), uint(42)).
		Return(errors.New("no tenés permiso para cancelar esta entrada"))

	r := setupTicketRouter(controllers.NewTicketController(svc), 42)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/tickets/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCancelTicketEndpoint_IDInvalido(t *testing.T) {
	svc := new(mockTicketService)

	r := setupTicketRouter(controllers.NewTicketController(svc), 42)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/tickets/abc", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTransferTicketEndpoint_Exitoso(t *testing.T) {
	svc := new(mockTicketService)
	svc.On("TransferTicket", uint(1), uint(42), domain.TransferTicketRequest{DestinoEmail: "otro@mail.com"}).
		Return(nil)

	r := setupTicketRouter(controllers.NewTicketController(svc), 42)
	body, _ := json.Marshal(map[string]string{"destino_email": "otro@mail.com"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/tickets/1/transfer", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTransferTicketEndpoint_BodyInvalido(t *testing.T) {
	svc := new(mockTicketService)

	r := setupTicketRouter(controllers.NewTicketController(svc), 42)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/tickets/1/transfer", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetMineEndpoint_ErrorServicio(t *testing.T) {
	svc := new(mockTicketService)
	svc.On("GetMyTickets", uint(42)).Return(nil, errors.New("db error"))

	r := setupTicketRouter(controllers.NewTicketController(svc), 42)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/tickets/mine", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestTransferTicketEndpoint_IDInvalido(t *testing.T) {
	svc := new(mockTicketService)

	r := setupTicketRouter(controllers.NewTicketController(svc), 42)
	body, _ := json.Marshal(map[string]string{"destino_email": "otro@mail.com"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/tickets/abc/transfer", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTransferTicketEndpoint_ErrorNegocio(t *testing.T) {
	svc := new(mockTicketService)
	svc.On("TransferTicket", uint(1), uint(42), domain.TransferTicketRequest{DestinoEmail: "otro@mail.com"}).
		Return(errors.New("no podés transferirte la entrada a vos mismo"))

	r := setupTicketRouter(controllers.NewTicketController(svc), 42)
	body, _ := json.Marshal(map[string]string{"destino_email": "otro@mail.com"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/tickets/1/transfer", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
