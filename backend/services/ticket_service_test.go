package services_test

import (
	"backend/domain"
	"backend/services"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type mockTicketDAO struct{ mock.Mock }

func (m *mockTicketDAO) Create(ticket *domain.Ticket) error {
	args := m.Called(ticket)
	return args.Error(0)
}
func (m *mockTicketDAO) FindByID(id uint) (*domain.Ticket, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Ticket), args.Error(1)
}
func (m *mockTicketDAO) FindByUserID(userID uint) ([]domain.Ticket, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Ticket), args.Error(1)
}
func (m *mockTicketDAO) Update(ticket *domain.Ticket) error {
	args := m.Called(ticket)
	return args.Error(0)
}

type mockWaitlistDAO struct{ mock.Mock }

func (m *mockWaitlistDAO) Create(entry *domain.WaitlistEntry) error {
	args := m.Called(entry)
	return args.Error(0)
}
func (m *mockWaitlistDAO) FindPendingByEventAndUser(eventID, userID uint) (*domain.WaitlistEntry, error) {
	args := m.Called(eventID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.WaitlistEntry), args.Error(1)
}
func (m *mockWaitlistDAO) FindFirstPending(eventID uint) (*domain.WaitlistEntry, error) {
	args := m.Called(eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.WaitlistEntry), args.Error(1)
}
func (m *mockWaitlistDAO) FindByUserID(userID uint) ([]domain.WaitlistEntry, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.WaitlistEntry), args.Error(1)
}
func (m *mockWaitlistDAO) Update(entry *domain.WaitlistEntry) error {
	args := m.Called(entry)
	return args.Error(0)
}

func activeTicket(ticketID, userID, eventID uint) *domain.Ticket {
	return &domain.Ticket{
		EventID:     eventID,
		UserID:      userID,
		Estado:      "activo",
		FechaCompra: time.Now(),
		Event: domain.Event{
			Titulo:         "Concierto",
			CapacidadTotal: 100,
		},
	}
}

// BuyTicket

func TestBuyTicket_Exitoso(t *testing.T) {
	td := new(mockTicketDAO)
	ed := new(mockEventDAO)
	ud := new(mockUserDAO)
	wd := new(mockWaitlistDAO)

	event := &domain.Event{CapacidadTotal: 10, EntradasVendidas: 5}
	ed.On("FindByID", uint(1)).Return(event, nil)
	td.On("Create", mock.AnythingOfType("*domain.Ticket")).Return(nil)
	ed.On("Update", mock.AnythingOfType("*domain.Event")).Return(nil)

	svc := services.NewTicketService(td, ed, ud, wd)
	resp, err := svc.BuyTicket(42, domain.BuyTicketRequest{EventID: 1})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "activo", resp.Estado)
	assert.Equal(t, 6, event.EntradasVendidas)
}

func TestBuyTicket_SinCupo(t *testing.T) {
	td := new(mockTicketDAO)
	ed := new(mockEventDAO)
	ud := new(mockUserDAO)
	wd := new(mockWaitlistDAO)

	event := &domain.Event{CapacidadTotal: 10, EntradasVendidas: 10}
	ed.On("FindByID", uint(1)).Return(event, nil)

	svc := services.NewTicketService(td, ed, ud, wd)
	resp, err := svc.BuyTicket(42, domain.BuyTicketRequest{EventID: 1})

	assert.Nil(t, resp)
	assert.EqualError(t, err, "no hay entradas disponibles")
}

func TestBuyTicket_EventoCancelado(t *testing.T) {
	td := new(mockTicketDAO)
	ed := new(mockEventDAO)
	ud := new(mockUserDAO)
	wd := new(mockWaitlistDAO)

	event := &domain.Event{CapacidadTotal: 10, Cancelado: true}
	ed.On("FindByID", uint(1)).Return(event, nil)

	svc := services.NewTicketService(td, ed, ud, wd)
	resp, err := svc.BuyTicket(42, domain.BuyTicketRequest{EventID: 1})

	assert.Nil(t, resp)
	assert.EqualError(t, err, "el evento está cancelado")
}

func TestBuyTicket_EventoNoExiste(t *testing.T) {
	td := new(mockTicketDAO)
	ed := new(mockEventDAO)
	ud := new(mockUserDAO)
	wd := new(mockWaitlistDAO)

	ed.On("FindByID", uint(99)).Return(nil, gorm.ErrRecordNotFound)

	svc := services.NewTicketService(td, ed, ud, wd)
	resp, err := svc.BuyTicket(42, domain.BuyTicketRequest{EventID: 99})

	assert.Nil(t, resp)
	assert.EqualError(t, err, "evento no encontrado")
}

// GetMyTickets

func TestGetMyTickets_Exitoso(t *testing.T) {
	td := new(mockTicketDAO)
	ed := new(mockEventDAO)
	ud := new(mockUserDAO)
	wd := new(mockWaitlistDAO)

	tickets := []domain.Ticket{*activeTicket(1, 42, 1)}
	td.On("FindByUserID", uint(42)).Return(tickets, nil)

	svc := services.NewTicketService(td, ed, ud, wd)
	result, err := svc.GetMyTickets(42)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "activo", result[0].Estado)
}

func TestGetMyTickets_ErrorDAO(t *testing.T) {
	td := new(mockTicketDAO)
	ed := new(mockEventDAO)
	ud := new(mockUserDAO)
	wd := new(mockWaitlistDAO)

	td.On("FindByUserID", uint(42)).Return(nil, errors.New("db error"))

	svc := services.NewTicketService(td, ed, ud, wd)
	result, err := svc.GetMyTickets(42)

	assert.Nil(t, result)
	assert.Error(t, err)
}

// CancelTicket

func TestCancelTicket_Exitoso(t *testing.T) {
	td := new(mockTicketDAO)
	ed := new(mockEventDAO)
	ud := new(mockUserDAO)
	wd := new(mockWaitlistDAO)

	ticket := activeTicket(1, 42, 1)
	event := &domain.Event{CapacidadTotal: 10, EntradasVendidas: 3}

	td.On("FindByID", uint(1)).Return(ticket, nil)
	ed.On("FindByID", uint(1)).Return(event, nil)
	ed.On("Update", mock.AnythingOfType("*domain.Event")).Return(nil)
	td.On("Update", mock.AnythingOfType("*domain.Ticket")).Return(nil)
	// Sin nadie en la lista de espera: el cupo se libera.
	wd.On("FindFirstPending", uint(1)).Return(nil, nil)

	svc := services.NewTicketService(td, ed, ud, wd)
	err := svc.CancelTicket(1, 42)

	assert.NoError(t, err)
	assert.Equal(t, "cancelado", ticket.Estado)
	assert.Equal(t, 2, event.EntradasVendidas)
}

func TestCancelTicket_AsignaAListaDeEspera(t *testing.T) {
	td := new(mockTicketDAO)
	ed := new(mockEventDAO)
	ud := new(mockUserDAO)
	wd := new(mockWaitlistDAO)

	ticket := activeTicket(1, 42, 1)
	event := &domain.Event{CapacidadTotal: 10, EntradasVendidas: 10}
	entry := &domain.WaitlistEntry{EventID: 1, UserID: 77, Estado: "pendiente"}

	td.On("FindByID", uint(1)).Return(ticket, nil)
	ed.On("FindByID", uint(1)).Return(event, nil)
	td.On("Update", mock.AnythingOfType("*domain.Ticket")).Return(nil)
	wd.On("FindFirstPending", uint(1)).Return(entry, nil)
	// Se le crea un ticket al primero de la lista.
	td.On("Create", mock.AnythingOfType("*domain.Ticket")).Return(nil)
	wd.On("Update", mock.AnythingOfType("*domain.WaitlistEntry")).Return(nil)

	svc := services.NewTicketService(td, ed, ud, wd)
	err := svc.CancelTicket(1, 42)

	assert.NoError(t, err)
	assert.Equal(t, "cancelado", ticket.Estado)
	// El cupo NO se libera: la butaca pasa al de la lista.
	assert.Equal(t, 10, event.EntradasVendidas)
	assert.Equal(t, "asignado", entry.Estado)
	assert.NotNil(t, entry.FechaAsignacion)
	// El evento no se actualiza porque el cupo no cambió.
	ed.AssertNotCalled(t, "Update", mock.Anything)
}

func TestCancelTicket_NoEsDueno(t *testing.T) {
	td := new(mockTicketDAO)
	ed := new(mockEventDAO)
	ud := new(mockUserDAO)
	wd := new(mockWaitlistDAO)

	ticket := activeTicket(1, 42, 1)
	td.On("FindByID", uint(1)).Return(ticket, nil)

	svc := services.NewTicketService(td, ed, ud, wd)
	err := svc.CancelTicket(1, 99)

	assert.EqualError(t, err, "no tenés permiso para cancelar esta entrada")
}

func TestCancelTicket_NoActiva(t *testing.T) {
	td := new(mockTicketDAO)
	ed := new(mockEventDAO)
	ud := new(mockUserDAO)
	wd := new(mockWaitlistDAO)

	ticket := activeTicket(1, 42, 1)
	ticket.Estado = "cancelado"
	td.On("FindByID", uint(1)).Return(ticket, nil)

	svc := services.NewTicketService(td, ed, ud, wd)
	err := svc.CancelTicket(1, 42)

	assert.EqualError(t, err, "la entrada no está activa")
}

// TransferTicket

func TestTransferTicket_Exitoso(t *testing.T) {
	td := new(mockTicketDAO)
	ed := new(mockEventDAO)
	ud := new(mockUserDAO)
	wd := new(mockWaitlistDAO)

	ticket := activeTicket(1, 42, 1)
	destUser := &domain.User{Email: "destino@mail.com"}
	destUser.ID = 99

	td.On("FindByID", uint(1)).Return(ticket, nil)
	ud.On("FindByEmail", "destino@mail.com").Return(destUser, nil)
	td.On("Update", mock.AnythingOfType("*domain.Ticket")).Return(nil)

	svc := services.NewTicketService(td, ed, ud, wd)
	err := svc.TransferTicket(1, 42, domain.TransferTicketRequest{DestinoEmail: "destino@mail.com"})

	assert.NoError(t, err)
	assert.Equal(t, uint(99), ticket.UserID)
}

func TestTransferTicket_NoEsDueno(t *testing.T) {
	td := new(mockTicketDAO)
	ed := new(mockEventDAO)
	ud := new(mockUserDAO)
	wd := new(mockWaitlistDAO)

	ticket := activeTicket(1, 42, 1)
	td.On("FindByID", uint(1)).Return(ticket, nil)

	svc := services.NewTicketService(td, ed, ud, wd)
	err := svc.TransferTicket(1, 99, domain.TransferTicketRequest{DestinoEmail: "destino@mail.com"})

	assert.EqualError(t, err, "no tenés permiso para transferir esta entrada")
}

func TestTransferTicket_DestinoNoExiste(t *testing.T) {
	td := new(mockTicketDAO)
	ed := new(mockEventDAO)
	ud := new(mockUserDAO)
	wd := new(mockWaitlistDAO)

	ticket := activeTicket(1, 42, 1)
	td.On("FindByID", uint(1)).Return(ticket, nil)
	ud.On("FindByEmail", "noexiste@mail.com").Return(nil, gorm.ErrRecordNotFound)

	svc := services.NewTicketService(td, ed, ud, wd)
	err := svc.TransferTicket(1, 42, domain.TransferTicketRequest{DestinoEmail: "noexiste@mail.com"})

	assert.EqualError(t, err, "usuario destino no encontrado")
}

func TestTransferTicket_TransferenciaASiMismo(t *testing.T) {
	td := new(mockTicketDAO)
	ed := new(mockEventDAO)
	ud := new(mockUserDAO)
	wd := new(mockWaitlistDAO)

	ticket := activeTicket(1, 42, 1)
	selfUser := &domain.User{Email: "self@mail.com"}
	selfUser.ID = 42

	td.On("FindByID", uint(1)).Return(ticket, nil)
	ud.On("FindByEmail", "self@mail.com").Return(selfUser, nil)

	svc := services.NewTicketService(td, ed, ud, wd)
	err := svc.TransferTicket(1, 42, domain.TransferTicketRequest{DestinoEmail: "self@mail.com"})

	assert.EqualError(t, err, "no podés transferirte la entrada a vos mismo")
}
