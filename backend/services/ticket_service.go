package services

import (
	"backend/dao"
	"backend/domain"
	"errors"
)

type ticketService struct {
	ticketDAO dao.ITicketDAO
	eventDAO  dao.IEventDAO
	userDAO   dao.IUserDAO
}

func NewTicketService(td dao.ITicketDAO, ed dao.IEventDAO, ud dao.IUserDAO) ITicketService {
	return &ticketService{ticketDAO: td, eventDAO: ed, userDAO: ud}
}

func (s *ticketService) BuyTicket(_ uint, _ domain.BuyTicketRequest) (*domain.TicketResponse, error) {
	return nil, errors.New("not implemented")
}

func (s *ticketService) GetMyTickets(_ uint) ([]domain.TicketResponse, error) {
	return nil, errors.New("not implemented")
}

func (s *ticketService) CancelTicket(_, _ uint) error {
	return errors.New("not implemented")
}

func (s *ticketService) TransferTicket(_ uint, _ uint, _ domain.TransferTicketRequest) error {
	return errors.New("not implemented")
}
