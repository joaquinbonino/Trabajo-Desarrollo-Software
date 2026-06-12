package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// dummyOK es un handler de prueba que responde 200 y refleja el contexto.
func dummyOK(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"user_id": c.GetUint("user_id"),
		"rol":     c.GetString("rol"),
	})
}

// AuthMiddleware

func TestAuthMiddleware_SinHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AuthMiddleware())
	r.GET("/protegido", dummyOK)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protegido", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_SinPrefijoBearer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AuthMiddleware())
	r.GET("/protegido", dummyOK)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protegido", nil)
	req.Header.Set("Authorization", "Token abc123")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_TokenInvalido(t *testing.T) {
	t.Setenv("JWT_SECRET", "clave-de-prueba")
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AuthMiddleware())
	r.GET("/protegido", dummyOK)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protegido", nil)
	req.Header.Set("Authorization", "Bearer token.invalido")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_TokenValido(t *testing.T) {
	t.Setenv("JWT_SECRET", "clave-de-prueba")
	token, err := GenerateToken(7, "admin")
	assert.NoError(t, err)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AuthMiddleware())
	r.GET("/protegido", dummyOK)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protegido", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "admin")
	assert.Contains(t, w.Body.String(), "7")
}

// AdminMiddleware

func TestAdminMiddleware_RolNoAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("rol", "cliente"); c.Next() })
	r.Use(AdminMiddleware())
	r.GET("/admin", dummyOK)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/admin", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestAdminMiddleware_RolAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("rol", "admin"); c.Next() })
	r.Use(AdminMiddleware())
	r.GET("/admin", dummyOK)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/admin", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// CORSMiddleware

func TestCORSMiddleware_OptionsDevuelve204(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(CORSMiddleware())
	r.GET("/x", dummyOK)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodOptions, "/x", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORSMiddleware_RequestNormal(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(CORSMiddleware())
	r.GET("/x", dummyOK)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/x", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
}
