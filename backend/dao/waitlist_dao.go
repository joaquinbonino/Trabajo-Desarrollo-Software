package dao

import (
	"backend/domain"
	"errors"

	"gorm.io/gorm"
)

type waitlistDAO struct {
	db *gorm.DB
}

func NewWaitlistDAO(db *gorm.DB) IWaitlistDAO {
	return &waitlistDAO{db: db}
}

func (d *waitlistDAO) Create(entry *domain.WaitlistEntry) error {
	return d.db.Create(entry).Error
}

func (d *waitlistDAO) FindPendingByEventAndUser(eventID, userID uint) (*domain.WaitlistEntry, error) {
	var entry domain.WaitlistEntry
	err := d.db.
		Where("event_id = ? AND user_id = ? AND estado = ?", eventID, userID, "pendiente").
		First(&entry).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (d *waitlistDAO) FindFirstPending(eventID uint) (*domain.WaitlistEntry, error) {
	var entry domain.WaitlistEntry
	err := d.db.
		Where("event_id = ? AND estado = ?", eventID, "pendiente").
		Order("created_at asc").
		First(&entry).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (d *waitlistDAO) FindByUserID(userID uint) ([]domain.WaitlistEntry, error) {
	var entries []domain.WaitlistEntry
	err := d.db.Preload("Event").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&entries).Error
	return entries, err
}

func (d *waitlistDAO) Update(entry *domain.WaitlistEntry) error {
	return d.db.Omit("User", "Event").Save(entry).Error
}
