package controllers

import (
	"backend/domain"
	"backend/services"
	"backend/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TicketController struct {
	service services.ITicketService
}

func NewTicketController(s services.ITicketService) *TicketController {
	return &TicketController{service: s}
}

func (h *TicketController) Buy(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req domain.BuyTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	ticket, err := h.service.BuyTicket(userID, req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusCreated, ticket)
}

func (h *TicketController) GetMine(c *gin.Context) {
	userID := c.GetUint("user_id")

	tickets, err := h.service.GetMyTickets(userID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, tickets)
}

func (h *TicketController) Cancel(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id inválido")
		return
	}

	if err := h.service.CancelTicket(uint(id), userID); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"message": "entrada cancelada"})
}

func (h *TicketController) Transfer(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id inválido")
		return
	}

	var req domain.TransferTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.TransferTicket(uint(id), userID, req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"message": "entrada transferida"})
}
