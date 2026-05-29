package services

import (
	"backend/dao"
	"backend/domain"
	"errors"
)

type eventService struct {
	dao       dao.IEventDAO
	ticketDAO dao.ITicketDAO
}

func NewEventService(d dao.IEventDAO, td dao.ITicketDAO) IEventService {
	return &eventService{dao: d, ticketDAO: td}
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

func (s *eventService) UpdateEvent(id uint, req domain.UpdateEventRequest) (*domain.EventResponse, error) {
	event, err := s.dao.FindByID(id)
	if err != nil {
		return nil, errors.New("evento no encontrado")
	}
	if event.Cancelado {
		return nil, errors.New("no se puede modificar un evento cancelado")
	}
	if req.Titulo != "" {
		event.Titulo = req.Titulo
	}
	if req.Descripcion != "" {
		event.Descripcion = req.Descripcion
	}
	if req.Categoria != "" {
		event.Categoria = req.Categoria
	}
	if !req.FechaHora.IsZero() {
		event.FechaHora = req.FechaHora
	}
	if req.Duracion != 0 {
		event.Duracion = req.Duracion
	}
	if req.CapacidadTotal != 0 {
		if req.CapacidadTotal < event.EntradasVendidas {
			return nil, errors.New("la capacidad no puede ser menor a las entradas ya vendidas")
		}
		event.CapacidadTotal = req.CapacidadTotal
	}
	if req.Foto != "" {
		event.Foto = req.Foto
	}
	if err := s.dao.Update(event); err != nil {
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

func (s *eventService) GetEventReport(id uint) (*domain.EventReportResponse, error) {
	event, err := s.dao.FindByID(id)
	if err != nil {
		return nil, errors.New("evento no encontrado")
	}
	tickets, err := s.ticketDAO.FindByEventID(id)
	if err != nil {
		return nil, err
	}
	buyers := make([]domain.BuyerInfo, 0, len(tickets))
	for _, t := range tickets {
		buyers = append(buyers, domain.BuyerInfo{
			UserID:      t.UserID,
			Nombre:      t.User.Nombre,
			Email:       t.User.Email,
			Estado:      t.Estado,
			FechaCompra: t.FechaCompra,
		})
	}
	return &domain.EventReportResponse{
		EventID:             event.ID,
		Titulo:              event.Titulo,
		CapacidadTotal:      event.CapacidadTotal,
		EntradasVendidas:    event.EntradasVendidas,
		EntradasDisponibles: event.CapacidadTotal - event.EntradasVendidas,
		Compradores:         buyers,
	}, nil
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
