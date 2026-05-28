package controllers

import (
	"backend/services"
	"backend/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EventController struct {
	service services.IEventService
}

func NewEventController(s services.IEventService) *EventController {
	return &EventController{service: s}
}

func (h *EventController) List(c *gin.Context) {
	categoria := c.Query("categoria")
	events, err := h.service.ListEvents(categoria)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, events)
}

func (h *EventController) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id inválido")
		return
	}
	event, err := h.service.GetEvent(uint(id))
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, event)
}
