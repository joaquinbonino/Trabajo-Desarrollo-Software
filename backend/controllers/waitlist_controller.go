package controllers

import (
	"backend/services"
	"backend/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WaitlistController struct {
	service services.IWaitlistService
}

func NewWaitlistController(s services.IWaitlistService) *WaitlistController {
	return &WaitlistController{service: s}
}

func (h *WaitlistController) Join(c *gin.Context) {
	userID := c.GetUint("user_id")

	eventID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id inválido")
		return
	}

	if err := h.service.JoinWaitlist(userID, uint(eventID)); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusCreated, gin.H{"message": "te anotaste en la lista de espera"})
}

func (h *WaitlistController) GetMine(c *gin.Context) {
	userID := c.GetUint("user_id")

	entries, err := h.service.GetMyWaitlist(userID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, entries)
}
