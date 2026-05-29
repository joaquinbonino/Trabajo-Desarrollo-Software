package controllers

import (
	"backend/domain"
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

func (h *EventController) Create(c *gin.Context) {
	var req domain.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	event, err := h.service.CreateEvent(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, http.StatusCreated, event)
}

func (h *EventController) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id inválido")
		return
	}
	var req domain.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	event, err := h.service.UpdateEvent(uint(id), req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, event)
}

func (h *EventController) Cancel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id inválido")
		return
	}
	if err := h.service.CancelEvent(uint(id)); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"message": "evento cancelado"})
}
