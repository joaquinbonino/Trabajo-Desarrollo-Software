package controllers_test

import (
	"backend/controllers"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// TestHealthEndpoint_DBNoDisponible cubre la rama de error del health check.
// Usa el driver MySQL (ya dependencia del proyecto, sin agregar nada nuevo) con
// un DSN apuntando a un puerto muerto: SkipInitializeWithVersion evita que
// gorm.Open intente conectar, y el `SELECT 1` falla al ejecutarse → 503.
func TestHealthEndpoint_DBNoDisponible(t *testing.T) {
	sqlDB, _ := sql.Open("mysql", "user:pass@tcp(127.0.0.1:1)/nodb")

	// gorm.Open puede devolver error al no poder conectar; eso es justamente la
	// condición "DB caída" que queremos. GORM igual devuelve un *gorm.DB usable
	// cuyo Exec fallará dentro del handler.
	db, _ := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	gin.SetMode(gin.TestMode)
	r := gin.New()
	ctrl := controllers.NewHealthController(db)
	r.GET("/health", ctrl.Check)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}
