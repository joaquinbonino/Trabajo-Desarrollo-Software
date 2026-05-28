package domain

import "time"

// Auth
type RegisterRequest struct {
	Nombre   string `json:"nombre" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID     uint   `json:"id"`
	Nombre string `json:"nombre"`
	Email  string `json:"email"`
	Rol    string `json:"rol"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// Events
type CreateEventRequest struct {
	Titulo         string    `json:"titulo" binding:"required"`
	Descripcion    string    `json:"descripcion"`
	Categoria      string    `json:"categoria"`
	FechaHora      time.Time `json:"fecha_hora" binding:"required"`
	Duracion       int       `json:"duracion"`
	CapacidadTotal int       `json:"capacidad_total" binding:"required,min=1"`
	Foto           string    `json:"foto"`
}

type EventResponse struct {
	ID               uint      `json:"id"`
	Titulo           string    `json:"titulo"`
	Descripcion      string    `json:"descripcion"`
	Categoria        string    `json:"categoria"`
	FechaHora        time.Time `json:"fecha_hora"`
	Duracion         int       `json:"duracion"`
	CapacidadTotal   int       `json:"capacidad_total"`
	EntradasVendidas int       `json:"entradas_vendidas"`
	Foto             string    `json:"foto"`
	Cancelado        bool      `json:"cancelado"`
}

// Tickets
type BuyTicketRequest struct {
	EventID uint `json:"event_id" binding:"required"`
}

type TransferTicketRequest struct {
	DestinoEmail string `json:"destino_email" binding:"required,email"`
}

type TicketResponse struct {
	ID          uint          `json:"id"`
	Estado      string        `json:"estado"`
	FechaCompra time.Time     `json:"fecha_compra"`
	Event       EventResponse `json:"event"`
}
