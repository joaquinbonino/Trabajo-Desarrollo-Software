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

// JoinWaitlist

func TestJoinWaitlist_Exitoso(t *testing.T) {
	wd := new(mockWaitlistDAO)
	ed := new(mockEventDAO)

	event := &domain.Event{CapacidadTotal: 10, EntradasVendidas: 10}
	ed.On("FindByID", uint(1)).Return(event, nil)
	wd.On("FindPendingByEventAndUser", uint(1), uint(42)).Return(nil, nil)
	wd.On("Create", mock.AnythingOfType("*domain.WaitlistEntry")).Return(nil)

	svc := services.NewWaitlistService(wd, ed)
	err := svc.JoinWaitlist(42, 1)

	assert.NoError(t, err)
	wd.AssertCalled(t, "Create", mock.AnythingOfType("*domain.WaitlistEntry"))
}

func TestJoinWaitlist_EventoNoExiste(t *testing.T) {
	wd := new(mockWaitlistDAO)
	ed := new(mockEventDAO)

	ed.On("FindByID", uint(99)).Return(nil, gorm.ErrRecordNotFound)

	svc := services.NewWaitlistService(wd, ed)
	err := svc.JoinWaitlist(42, 99)

	assert.EqualError(t, err, "evento no encontrado")
}

func TestJoinWaitlist_EventoCancelado(t *testing.T) {
	wd := new(mockWaitlistDAO)
	ed := new(mockEventDAO)

	event := &domain.Event{CapacidadTotal: 10, EntradasVendidas: 10, Cancelado: true}
	ed.On("FindByID", uint(1)).Return(event, nil)

	svc := services.NewWaitlistService(wd, ed)
	err := svc.JoinWaitlist(42, 1)

	assert.EqualError(t, err, "el evento está cancelado")
}

func TestJoinWaitlist_HayCupo(t *testing.T) {
	wd := new(mockWaitlistDAO)
	ed := new(mockEventDAO)

	event := &domain.Event{CapacidadTotal: 10, EntradasVendidas: 5}
	ed.On("FindByID", uint(1)).Return(event, nil)

	svc := services.NewWaitlistService(wd, ed)
	err := svc.JoinWaitlist(42, 1)

	assert.EqualError(t, err, "todavía hay entradas disponibles, comprá directamente")
}

func TestJoinWaitlist_YaAnotado(t *testing.T) {
	wd := new(mockWaitlistDAO)
	ed := new(mockEventDAO)

	event := &domain.Event{CapacidadTotal: 10, EntradasVendidas: 10}
	existente := &domain.WaitlistEntry{EventID: 1, UserID: 42, Estado: "pendiente"}
	ed.On("FindByID", uint(1)).Return(event, nil)
	wd.On("FindPendingByEventAndUser", uint(1), uint(42)).Return(existente, nil)

	svc := services.NewWaitlistService(wd, ed)
	err := svc.JoinWaitlist(42, 1)

	assert.EqualError(t, err, "ya estás en la lista de espera de este evento")
}

// GetMyWaitlist

func TestGetMyWaitlist_Exitoso(t *testing.T) {
	wd := new(mockWaitlistDAO)
	ed := new(mockEventDAO)

	ahora := time.Now()
	entries := []domain.WaitlistEntry{
		{Estado: "pendiente", Event: domain.Event{Titulo: "Recital"}},
		{Estado: "asignado", FechaAsignacion: &ahora, Event: domain.Event{Titulo: "Obra"}},
	}
	wd.On("FindByUserID", uint(42)).Return(entries, nil)

	svc := services.NewWaitlistService(wd, ed)
	result, err := svc.GetMyWaitlist(42)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "pendiente", result[0].Estado)
	assert.Equal(t, "asignado", result[1].Estado)
	assert.NotNil(t, result[1].FechaAsignacion)
}

func TestGetMyWaitlist_ErrorDAO(t *testing.T) {
	wd := new(mockWaitlistDAO)
	ed := new(mockEventDAO)

	wd.On("FindByUserID", uint(42)).Return(nil, errors.New("db error"))

	svc := services.NewWaitlistService(wd, ed)
	result, err := svc.GetMyWaitlist(42)

	assert.Nil(t, result)
	assert.Error(t, err)
}
