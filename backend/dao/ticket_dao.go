package dao

import (
	"backend/domain"

	"gorm.io/gorm"
)

type ticketDAO struct {
	db *gorm.DB
}

func NewTicketDAO(db *gorm.DB) ITicketDAO {
	return &ticketDAO{db: db}
}

func (d *ticketDAO) Create(ticket *domain.Ticket) error {
	return d.db.Create(ticket).Error
}

func (d *ticketDAO) FindByID(id uint) (*domain.Ticket, error) {
	var ticket domain.Ticket
	err := d.db.Preload("Event").Preload("User").First(&ticket, id).Error
	return &ticket, err
}

func (d *ticketDAO) FindByUserID(userID uint) ([]domain.Ticket, error) {
	var tickets []domain.Ticket
	err := d.db.Preload("Event").Where("user_id = ?", userID).Find(&tickets).Error
	return tickets, err
}

func (d *ticketDAO) Update(ticket *domain.Ticket) error {
	return d.db.Save(ticket).Error
}
