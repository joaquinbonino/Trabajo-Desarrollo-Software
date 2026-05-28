package dao

import (
	"backend/domain"

	"gorm.io/gorm"
)

type eventDAO struct {
	db *gorm.DB
}

func NewEventDAO(db *gorm.DB) IEventDAO {
	return &eventDAO{db: db}
}

func (d *eventDAO) Create(event *domain.Event) error {
	return d.db.Create(event).Error
}

func (d *eventDAO) FindByID(id uint) (*domain.Event, error) {
	var event domain.Event
	err := d.db.First(&event, id).Error
	return &event, err
}

func (d *eventDAO) FindAll(categoria string) ([]domain.Event, error) {
	var events []domain.Event
	q := d.db.Where("cancelado = false")
	if categoria != "" {
		q = q.Where("categoria = ?", categoria)
	}
	err := q.Find(&events).Error
	return events, err
}

func (d *eventDAO) Update(event *domain.Event) error {
	return d.db.Save(event).Error
}
