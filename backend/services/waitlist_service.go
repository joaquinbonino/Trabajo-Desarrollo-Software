package services

import (
	"backend/dao"
	"backend/domain"
	"errors"
)

type waitlistService struct {
	waitlistDAO dao.IWaitlistDAO
	eventDAO    dao.IEventDAO
}

func NewWaitlistService(wd dao.IWaitlistDAO, ed dao.IEventDAO) IWaitlistService {
	return &waitlistService{waitlistDAO: wd, eventDAO: ed}
}

// JoinWaitlist anota al usuario en la lista de espera de un evento. Solo se
// permite cuando el evento está agotado y el usuario no está ya anotado.
func (s *waitlistService) JoinWaitlist(userID, eventID uint) error {
	event, err := s.eventDAO.FindByID(eventID)
	if err != nil {
		return errors.New("evento no encontrado")
	}
	if event.Cancelado {
		return errors.New("el evento está cancelado")
	}
	if event.EntradasVendidas < event.CapacidadTotal {
		return errors.New("todavía hay entradas disponibles, comprá directamente")
	}

	existente, err := s.waitlistDAO.FindPendingByEventAndUser(eventID, userID)
	if err != nil {
		return err
	}
	if existente != nil {
		return errors.New("ya estás en la lista de espera de este evento")
	}

	entry := &domain.WaitlistEntry{
		EventID: eventID,
		UserID:  userID,
		Estado:  "pendiente",
	}
	return s.waitlistDAO.Create(entry)
}

// GetMyWaitlist devuelve las anotaciones del usuario autenticado, incluyendo las
// que ya le fueron asignadas (registro de la notificación).
func (s *waitlistService) GetMyWaitlist(userID uint) ([]domain.WaitlistResponse, error) {
	entries, err := s.waitlistDAO.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	result := make([]domain.WaitlistResponse, len(entries))
	for i, e := range entries {
		result[i] = domain.WaitlistResponse{
			ID:              e.ID,
			Estado:          e.Estado,
			FechaAsignacion: e.FechaAsignacion,
			Event:           toEventResponse(e.Event),
		}
	}
	return result, nil
}
