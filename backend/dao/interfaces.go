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

type IWaitlistDAO interface {
	Create(entry *domain.WaitlistEntry) error
	// FindPendingByEventAndUser devuelve la anotación pendiente del usuario en ese
	// evento, o (nil, nil) si no está anotado.
	FindPendingByEventAndUser(eventID, userID uint) (*domain.WaitlistEntry, error)
	// FindFirstPending devuelve el primero de la lista (FIFO) en estado pendiente,
	// o (nil, nil) si la lista está vacía.
	FindFirstPending(eventID uint) (*domain.WaitlistEntry, error)
	FindByUserID(userID uint) ([]domain.WaitlistEntry, error)
	Update(entry *domain.WaitlistEntry) error
}
