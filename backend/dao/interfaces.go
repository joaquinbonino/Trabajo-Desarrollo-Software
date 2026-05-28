package dao

import "backend/domain"

type IUserDAO interface {
	Create(user *domain.User) error
	FindByID(id uint) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	Update(user *domain.User) error
}

type IEventDAO interface {
	Create(event *domain.Event) error
	FindByID(id uint) (*domain.Event, error)
	FindAll(categoria string) ([]domain.Event, error)
	Update(event *domain.Event) error
}

type ITicketDAO interface {
	Create(ticket *domain.Ticket) error
	FindByID(id uint) (*domain.Ticket, error)
	FindByUserID(userID uint) ([]domain.Ticket, error)
	Update(ticket *domain.Ticket) error
}
