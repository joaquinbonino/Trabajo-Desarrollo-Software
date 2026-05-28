package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Nombre       string `gorm:"type:varchar(100);not null"`
	Email        string `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash string `gorm:"type:varchar(64);not null"`
	PasswordSalt string `gorm:"type:varchar(32);not null"`
	Rol          string `gorm:"type:enum('cliente','admin');default:'cliente'"`
}

type Event struct {
	gorm.Model
	Titulo           string    `gorm:"not null"`
	Descripcion      string
	Categoria        string
	FechaHora        time.Time `gorm:"not null"`
	Duracion         int
	CapacidadTotal   int  `gorm:"not null"`
	EntradasVendidas int  `gorm:"default:0"`
	Foto             string
	Cancelado        bool `gorm:"default:false"`
}

type Ticket struct {
	gorm.Model
	EventID     uint      `gorm:"not null"`
	UserID      uint      `gorm:"not null"`
	Estado      string    `gorm:"type:enum('activo','cancelado','transferido');default:'activo'"`
	FechaCompra time.Time `gorm:"not null"`
	Event       Event     `gorm:"foreignKey:EventID;references:ID"`
	User        User      `gorm:"foreignKey:UserID;references:ID"`
}
