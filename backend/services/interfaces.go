package services

import "backend/domain"

type IUserService interface {
	Register(req domain.RegisterRequest) (*domain.AuthResponse, error)
	Login(req domain.LoginRequest) (*domain.AuthResponse, error)
}

type IEventService interface {
	ListEvents(categoria string) ([]domain.EventResponse, error)
	ListAllEvents() ([]domain.EventResponse, error)
	GetEvent(id uint) (*domain.EventResponse, error)
	CreateEvent(req domain.CreateEventRequest) (*domain.EventResponse, error)
	UpdateEvent(id uint, req domain.UpdateEventRequest) (*domain.EventResponse, error)
	CancelEvent(id uint) error
	GetEventReport(id uint) (*domain.EventReportResponse, error)
}

type ITicketService interface {
	BuyTicket(userID uint, req domain.BuyTicketRequest) (*domain.TicketResponse, error)
	GetMyTickets(userID uint) ([]domain.TicketResponse, error)
	CancelTicket(ticketID, userID uint) error
	TransferTicket(ticketID, userID uint, req domain.TransferTicketRequest) error
}

type IWaitlistService interface {
	JoinWaitlist(userID, eventID uint) error
	GetMyWaitlist(userID uint) ([]domain.WaitlistResponse, error)
}
