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

func (s *eventService) ListEvents(_ string) ([]domain.EventResponse, error) {
	return nil, errors.New("not implemented")
}

func (s *eventService) GetEvent(_ uint) (*domain.EventResponse, error) {
	return nil, errors.New("not implemented")
}

func (s *eventService) CreateEvent(_ domain.CreateEventRequest) (*domain.EventResponse, error) {
	return nil, errors.New("not implemented")
}

func (s *eventService) CancelEvent(_ uint) error {
	return errors.New("not implemented")
}
