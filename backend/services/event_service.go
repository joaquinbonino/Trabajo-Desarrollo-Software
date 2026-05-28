package services

import (
	"backend/dao"
	"backend/domain"
	"errors"
)

type eventService struct {
	dao dao.IEventDAO
}

func NewEventService(d dao.IEventDAO) IEventService {
	return &eventService{dao: d}
}

func (s *eventService) ListEvents(categoria string) ([]domain.EventResponse, error) {
	events, err := s.dao.FindAll(categoria)
	if err != nil {
		return nil, err
	}
	result := make([]domain.EventResponse, len(events))
	for i, e := range events {
		result[i] = toEventResponse(e)
	}
	return result, nil
}

func (s *eventService) GetEvent(id uint) (*domain.EventResponse, error) {
	event, err := s.dao.FindByID(id)
	if err != nil {
		return nil, errors.New("evento no encontrado")
	}
	resp := toEventResponse(*event)
	return &resp, nil
}

func (s *eventService) CreateEvent(req domain.CreateEventRequest) (*domain.EventResponse, error) {
	event := &domain.Event{
		Titulo:         req.Titulo,
		Descripcion:    req.Descripcion,
		Categoria:      req.Categoria,
		FechaHora:      req.FechaHora,
		Duracion:       req.Duracion,
		CapacidadTotal: req.CapacidadTotal,
		Foto:           req.Foto,
	}
	if err := s.dao.Create(event); err != nil {
		return nil, err
	}
	resp := toEventResponse(*event)
	return &resp, nil
}

func (s *eventService) CancelEvent(id uint) error {
	event, err := s.dao.FindByID(id)
	if err != nil {
		return errors.New("evento no encontrado")
	}
	if event.Cancelado {
		return errors.New("el evento ya está cancelado")
	}
	event.Cancelado = true
	return s.dao.Update(event)
}

func toEventResponse(e domain.Event) domain.EventResponse {
	return domain.EventResponse{
		ID:               e.ID,
		Titulo:           e.Titulo,
		Descripcion:      e.Descripcion,
		Categoria:        e.Categoria,
		FechaHora:        e.FechaHora,
		Duracion:         e.Duracion,
		CapacidadTotal:   e.CapacidadTotal,
		EntradasVendidas: e.EntradasVendidas,
		Foto:             e.Foto,
		Cancelado:        e.Cancelado,
	}
}
