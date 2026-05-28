package controllers

import (
	"backend/domain"
	"backend/services"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service services.IUserService
}

func NewAuthController(s services.IUserService) *AuthController {
	return &AuthController{service: s}
}

func (h *AuthController) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.service.Register(req)
	if err != nil {
		utils.Error(c, http.StatusConflict, err.Error())
		return
	}

	utils.Success(c, http.StatusCreated, resp)
}

func (h *AuthController) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.service.Login(req)
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, resp)
}
