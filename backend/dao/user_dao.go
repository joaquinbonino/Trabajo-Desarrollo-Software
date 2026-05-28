package dao

import (
	"backend/domain"

	"gorm.io/gorm"
)

type userDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) IUserDAO {
	return &userDAO{db: db}
}

func (d *userDAO) Create(user *domain.User) error {
	return d.db.Create(user).Error
}

func (d *userDAO) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	err := d.db.First(&user, id).Error
	return &user, err
}

func (d *userDAO) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := d.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (d *userDAO) Update(user *domain.User) error {
	return d.db.Save(user).Error
}
