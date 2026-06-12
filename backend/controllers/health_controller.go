package controllers

import (
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthController struct {
	db *gorm.DB
}

func NewHealthController(db *gorm.DB) *HealthController {
	return &HealthController{db: db}
}

func (h *HealthController) Check(c *gin.Context) {
	if err := h.db.Exec("SELECT 1").Error; err != nil {
		utils.Error(c, http.StatusServiceUnavailable, "db no disponible")
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"status": "ok", "db": "connected"})
}
