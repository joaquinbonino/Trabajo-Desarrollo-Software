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

type mockTicketDAOForEvent struct{ mock.Mock }

func (m *mockTicketDAOForEvent) Create(ticket *domain.Ticket) error {
	return m.Called(ticket).Error(0)
}
func (m *mockTicketDAOForEvent) FindByID(id uint) (*domain.Ticket, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Ticket), args.Error(1)
}
func (m *mockTicketDAOForEvent) FindByUserID(userID uint) ([]domain.Ticket, error) {
	args := m.Called(userID)
	return args.Get(0).([]domain.Ticket), args.Error(1)
}
func (m *mockTicketDAOForEvent) FindByEventID(eventID uint) ([]domain.Ticket, error) {
	args := m.Called(eventID)
	return args.Get(0).([]domain.Ticket), args.Error(1)
}
func (m *mockTicketDAOForEvent) Update(ticket *domain.Ticket) error {
	return m.Called(ticket).Error(0)
}

type mockEventDAO struct{ mock.Mock }

func (m *mockEventDAO) Create(event *domain.Event) error {
	args := m.Called(event)
	return args.Error(0)
}
func (m *mockEventDAO) FindByID(id uint) (*domain.Event, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Event), args.Error(1)
}
func (m *mockEventDAO) FindAll(categoria string) ([]domain.Event, error) {
	args := m.Called(categoria)
	return args.Get(0).([]domain.Event), args.Error(1)
}
func (m *mockEventDAO) FindAllAdmin() ([]domain.Event, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Event), args.Error(1)
}
func (m *mockEventDAO) Update(event *domain.Event) error {
	args := m.Called(event)
	return args.Error(0)
}

func sampleEvent() *domain.Event {
	return &domain.Event{
		Titulo:         "Concierto",
		Categoria:      "Música",
		FechaHora:      time.Now().Add(24 * time.Hour),
		CapacidadTotal: 100,
	}
}

func TestListEvents_Exitoso(t *testing.T) {
	d := new(mockEventDAO)
	d.On("FindAll", "").Return([]domain.Event{*sampleEvent()}, nil)

	svc := services.NewEventService(d, new(mockTicketDAOForEvent))
	result, err := svc.ListEvents("")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Concierto", result[0].Titulo)
}

func TestListEvents_ConFiltroCategoria(t *testing.T) {
	d := new(mockEventDAO)
	d.On("FindAll", "Música").Return([]domain.Event{*sampleEvent()}, nil)

	svc := services.NewEventService(d, new(mockTicketDAOForEvent))
	result, err := svc.ListEvents("Música")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Música", result[0].Categoria)
}

func TestListEvents_ErrorDAO(t *testing.T) {
	d := new(mockEventDAO)
	d.On("FindAll", "").Return([]domain.Event{}, errors.New("db error"))

	svc := services.NewEventService(d, new(mockTicketDAOForEvent))
	result, err := svc.ListEvents("")

	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestGetEvent_Exitoso(t *testing.T) {
	d := new(mockEventDAO)
	d.On("FindByID", uint(1)).Return(sampleEvent(), nil)

	svc := services.NewEventService(d, new(mockTicketDAOForEvent))
	result, err := svc.GetEvent(1)

	assert.NoError(t, err)
	assert.Equal(t, "Concierto", result.Titulo)
}

func TestGetEvent_NoExiste(t *testing.T) {
	d := new(mockEventDAO)
	d.On("FindByID", uint(99)).Return(nil, gorm.ErrRecordNotFound)

	svc := services.NewEventService(d, new(mockTicketDAOForEvent))
	result, err := svc.GetEvent(99)

	assert.Nil(t, result)
	assert.EqualError(t, err, "evento no encontrado")
}

func TestCreateEvent_Exitoso(t *testing.T) {
	d := new(mockEventDAO)
	d.On("Create", mock.AnythingOfType("*domain.Event")).Return(nil)

	svc := services.NewEventService(d, new(mockTicketDAOForEvent))
	result, err := svc.CreateEvent(domain.CreateEventRequest{
		Titulo:         "Festival",
		CapacidadTotal: 200,
		FechaHora:      time.Now().Add(48 * time.Hour),
	})

	assert.NoError(t, err)
	assert.Equal(t, "Festival", result.Titulo)
}

func TestCancelEvent_Exitoso(t *testing.T) {
	event := sampleEvent()
	d := new(mockEventDAO)
	d.On("FindByID", uint(1)).Return(event, nil)
	d.On("Update", mock.AnythingOfType("*domain.Event")).Return(nil)

	svc := services.NewEventService(d, new(mockTicketDAOForEvent))
	err := svc.CancelEvent(1)

	assert.NoError(t, err)
	assert.True(t, event.Cancelado)
}

func TestCancelEvent_YaCancelado(t *testing.T) {
	event := sampleEvent()
	event.Cancelado = true
	d := new(mockEventDAO)
	d.On("FindByID", uint(1)).Return(event, nil)

	svc := services.NewEventService(d, new(mockTicketDAOForEvent))
	err := svc.CancelEvent(1)

	assert.EqualError(t, err, "el evento ya está cancelado")
}

func TestCancelEvent_NoExiste(t *testing.T) {
	d := new(mockEventDAO)
	d.On("FindByID", uint(99)).Return(nil, gorm.ErrRecordNotFound)

	svc := services.NewEventService(d, new(mockTicketDAOForEvent))
	err := svc.CancelEvent(99)

	assert.EqualError(t, err, "evento no encontrado")
}

func TestCreateEvent_ErrorDAO(t *testing.T) {
	d := new(mockEventDAO)
	d.On("Create", mock.AnythingOfType("*domain.Event")).Return(errors.New("db error"))

	svc := services.NewEventService(d, new(mockTicketDAOForEvent))
	result, err := svc.CreateEvent(domain.CreateEventRequest{Titulo: "Festival", CapacidadTotal: 10})

	assert.Nil(t, result)
	assert.Error(t, err)
}

// UpdateEvent

func TestUpdateEvent_Exitoso(t *testing.T) {
	event := sampleEvent()
	event.EntradasVendidas = 10
	d := new(mockEventDAO)
	d.On("FindByID", uint(1)).Return(event, nil)
	d.On("Update", mock.AnythingOfType("*domain.Event")).Return(nil)

	svc := services.NewEventService(d, new(mockTicketDAOForEvent))
	result, err := svc.UpdateEvent(1, domain.UpdateEventRequest{
		Titulo:         "Nuevo título",
		Descripcion:    "Nueva descripción",
		Categoria:      "Teatro",
		Duracion:       120,
		CapacidadTotal: 150,
		Foto:           "foto.jpg",
		FechaHora:      time.Now().Add(72 * time.Hour),
	})

	assert.NoError(t, err)
	assert.Equal(t, "Nuevo título", result.Titulo)
	assert.Equal(t, "Teatro", result.Categoria)
	assert.Equal(t, 150, result.CapacidadTotal)
}

func TestUpdateEvent_NoExiste(t *testing.T) {
	d := new(mockEventDAO)
	d.On("FindByID", uint(99)).Return(nil, gorm.ErrRecordNotFound)

	svc := services.NewEventService(d, new(mockTicketDAOForEvent))
	result, err := svc.UpdateEvent(99, domain.UpdateEventRequest{Titulo: "x"})

	assert.Nil(t, result)
	assert.EqualError(t, err, "evento no encontrado")
}

func TestUpdateEvent_EventoCancelado(t *testing.T) {
	event := sampleEvent()
	event.Cancelado = true
	d := new(mockEventDAO)
	d.On("FindByID", uint(1)).Return(event, nil)

	svc := services.NewEventService(d, new(mockTicketDAOForEvent))
	result, err := svc.UpdateEvent(1, domain.UpdateEventRequest{Titulo: "x"})

	assert.Nil(t, result)
	assert.EqualError(t, err, "no se puede modificar un evento cancelado")
}

func TestUpdateEvent_CapacidadMenorAVendidas(t *testing.T) {
	event := sampleEvent()
	event.EntradasVendidas = 50
	d := new(mockEventDAO)
	d.On("FindByID", uint(1)).Return(event, nil)

	svc := services.NewEventService(d, new(mockTicketDAOForEvent))
	result, err := svc.UpdateEvent(1, domain.UpdateEventRequest{CapacidadTotal: 10})

	assert.Nil(t, result)
	assert.EqualError(t, err, "la capacidad no puede ser menor a las entradas ya vendidas")
}

func TestUpdateEvent_ErrorDAO(t *testing.T) {
	event := sampleEvent()
	d := new(mockEventDAO)
	d.On("FindByID", uint(1)).Return(event, nil)
	d.On("Update", mock.AnythingOfType("*domain.Event")).Return(errors.New("db error"))

	svc := services.NewEventService(d, new(mockTicketDAOForEvent))
	result, err := svc.UpdateEvent(1, domain.UpdateEventRequest{Titulo: "x"})

	assert.Nil(t, result)
	assert.Error(t, err)
}

// ListAllEvents

func TestListAllEvents_Exitoso(t *testing.T) {
	cancelado := sampleEvent()
	cancelado.Cancelado = true
	d := new(mockEventDAO)
	d.On("FindAllAdmin").Return([]domain.Event{*sampleEvent(), *cancelado}, nil)

	svc := services.NewEventService(d, new(mockTicketDAOForEvent))
	result, err := svc.ListAllEvents()

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestListAllEvents_ErrorDAO(t *testing.T) {
	d := new(mockEventDAO)
	d.On("FindAllAdmin").Return(nil, errors.New("db error"))

	svc := services.NewEventService(d, new(mockTicketDAOForEvent))
	result, err := svc.ListAllEvents()

	assert.Nil(t, result)
	assert.Error(t, err)
}

// GetEventReport

func TestGetEventReport_Exitoso(t *testing.T) {
	event := sampleEvent()
	event.EntradasVendidas = 2
	ed := new(mockEventDAO)
	td := new(mockTicketDAOForEvent)
	ed.On("FindByID", uint(1)).Return(event, nil)
	tickets := []domain.Ticket{
		{UserID: 10, Estado: "activo", FechaCompra: time.Now(), User: domain.User{Nombre: "Ana", Email: "ana@mail.com"}},
		{UserID: 11, Estado: "cancelado", FechaCompra: time.Now(), User: domain.User{Nombre: "Beto", Email: "beto@mail.com"}},
	}
	td.On("FindByEventID", uint(1)).Return(tickets, nil)

	svc := services.NewEventService(ed, td)
	report, err := svc.GetEventReport(1)

	assert.NoError(t, err)
	assert.Equal(t, 100, report.CapacidadTotal)
	assert.Equal(t, 2, report.EntradasVendidas)
	assert.Equal(t, 98, report.EntradasDisponibles)
	assert.Len(t, report.Compradores, 2)
	assert.Equal(t, "Ana", report.Compradores[0].Nombre)
}

func TestGetEventReport_NoExiste(t *testing.T) {
	ed := new(mockEventDAO)
	ed.On("FindByID", uint(99)).Return(nil, gorm.ErrRecordNotFound)

	svc := services.NewEventService(ed, new(mockTicketDAOForEvent))
	report, err := svc.GetEventReport(99)

	assert.Nil(t, report)
	assert.EqualError(t, err, "evento no encontrado")
}

func TestGetEventReport_ErrorTicketDAO(t *testing.T) {
	event := sampleEvent()
	ed := new(mockEventDAO)
	td := new(mockTicketDAOForEvent)
	ed.On("FindByID", uint(1)).Return(event, nil)
	td.On("FindByEventID", uint(1)).Return([]domain.Ticket{}, errors.New("db error"))

	svc := services.NewEventService(ed, td)
	report, err := svc.GetEventReport(1)

	assert.Nil(t, report)
	assert.Error(t, err)
}
