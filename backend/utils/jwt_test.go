package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateYValidateToken_RoundTrip(t *testing.T) {
	t.Setenv("JWT_SECRET", "clave-de-prueba")

	token, err := GenerateToken(7, "admin")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, uint(7), claims.UserID)
	assert.Equal(t, "admin", claims.Rol)
}

func TestValidateToken_Malformado(t *testing.T) {
	t.Setenv("JWT_SECRET", "clave-de-prueba")

	claims, err := ValidateToken("esto.no.es.un.token")
	assert.Nil(t, claims)
	assert.Error(t, err)
}

func TestValidateToken_SecretoIncorrecto(t *testing.T) {
	t.Setenv("JWT_SECRET", "secreto-original")
	token, err := GenerateToken(1, "cliente")
	assert.NoError(t, err)

	// El token se valida con otro secreto: la firma no coincide.
	t.Setenv("JWT_SECRET", "otro-secreto")
	claims, err := ValidateToken(token)
	assert.Nil(t, claims)
	assert.Error(t, err)
}

func TestValidateToken_Expirado(t *testing.T) {
	t.Setenv("JWT_SECRET", "clave-de-prueba")

	claims := Claims{
		UserID: 1,
		Rol:    "cliente",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("clave-de-prueba"))
	assert.NoError(t, err)

	parsed, err := ValidateToken(token)
	assert.Nil(t, parsed)
	assert.Error(t, err)
}

func TestValidateToken_MetodoDeFirmaInesperado(t *testing.T) {
	t.Setenv("JWT_SECRET", "clave-de-prueba")

	// Token firmado con "none" (algoritmo distinto a HMAC): debe rechazarse.
	claims := Claims{UserID: 1, Rol: "cliente"}
	token, err := jwt.NewWithClaims(jwt.SigningMethodNone, claims).SignedString(jwt.UnsafeAllowNoneSignatureType)
	assert.NoError(t, err)

	parsed, err := ValidateToken(token)
	assert.Nil(t, parsed)
	assert.Error(t, err)
}
