package services

import (
	"backend/dao"
	"backend/domain"
	"errors"
	"time"
)

type ticketService struct {
	ticketDAO dao.ITicketDAO
	eventDAO  dao.IEventDAO
	userDAO   dao.IUserDAO
}

func NewTicketService(td dao.ITicketDAO, ed dao.IEventDAO, ud dao.IUserDAO) ITicketService {
	return &ticketService{ticketDAO: td, eventDAO: ed, userDAO: ud}
}

func (s *ticketService) BuyTicket(userID uint, req domain.BuyTicketRequest) (*domain.TicketResponse, error) {
	event, err := s.eventDAO.FindByID(req.EventID)
	if err != nil {
		return nil, errors.New("evento no encontrado")
	}
	if event.Cancelado {
		return nil, errors.New("el evento está cancelado")
	}
	if event.EntradasVendidas >= event.CapacidadTotal {
		return nil, errors.New("no hay entradas disponibles")
	}

	ticket := &domain.Ticket{
		EventID:     req.EventID,
		UserID:      userID,
		Estado:      "activo",
		FechaCompra: time.Now(),
	}
	if err := s.ticketDAO.Create(ticket); err != nil {
		return nil, err
	}

	event.EntradasVendidas++
	if err := s.eventDAO.Update(event); err != nil {
		return nil, err
	}

	resp := &domain.TicketResponse{
		ID:          ticket.ID,
		Estado:      ticket.Estado,
		FechaCompra: ticket.FechaCompra,
		Event:       toEventResponse(*event),
	}
	return resp, nil
}

func (s *ticketService) GetMyTickets(userID uint) ([]domain.TicketResponse, error) {
	tickets, err := s.ticketDAO.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	result := make([]domain.TicketResponse, len(tickets))
	for i, t := range tickets {
		result[i] = domain.TicketResponse{
			ID:          t.ID,
			Estado:      t.Estado,
			FechaCompra: t.FechaCompra,
			Event:       toEventResponse(t.Event),
		}
	}
	return result, nil
}

func (s *ticketService) CancelTicket(ticketID, userID uint) error {
	ticket, err := s.ticketDAO.FindByID(ticketID)
	if err != nil {
		return errors.New("entrada no encontrada")
	}
	if ticket.UserID != userID {
		return errors.New("no tenés permiso para cancelar esta entrada")
	}
	if ticket.Estado != "activo" {
		return errors.New("la entrada no está activa")
	}

	event, err := s.eventDAO.FindByID(ticket.EventID)
	if err != nil {
		return errors.New("evento no encontrado")
	}
	if event.EntradasVendidas > 0 {
		event.EntradasVendidas--
	}
	if err := s.eventDAO.Update(event); err != nil {
		return err
	}

	ticket.Estado = "cancelado"
	return s.ticketDAO.Update(ticket)
}

func (s *ticketService) TransferTicket(ticketID, userID uint, req domain.TransferTicketRequest) error {
	ticket, err := s.ticketDAO.FindByID(ticketID)
	if err != nil {
		return errors.New("entrada no encontrada")
	}
	if ticket.UserID != userID {
		return errors.New("no tenés permiso para transferir esta entrada")
	}
	if ticket.Estado != "activo" {
		return errors.New("la entrada no está activa")
	}

	destUser, err := s.userDAO.FindByEmail(req.DestinoEmail)
	if err != nil {
		return errors.New("usuario destino no encontrado")
	}
	if destUser.ID == userID {
		return errors.New("no podés transferirte la entrada a vos mismo")
	}

	ticket.UserID = destUser.ID
	return s.ticketDAO.Update(ticket)
}
