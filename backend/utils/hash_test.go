package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword_Deterministico(t *testing.T) {
	h1 := HashPassword("secret", "salt123")
	h2 := HashPassword("secret", "salt123")
	assert.Equal(t, h1, h2)
}

func TestHashPassword_DistintoSalt(t *testing.T) {
	h1 := HashPassword("secret", "salt1")
	h2 := HashPassword("secret", "salt2")
	assert.NotEqual(t, h1, h2)
}

func TestHashPassword_NoExponeLaPassword(t *testing.T) {
	h := HashPassword("miClaveSuperSecreta", "salt")
	assert.NotContains(t, h, "miClaveSuperSecreta")
	assert.Len(t, h, 64) // SHA-256 en hex son 64 caracteres
}

func TestGenerateSalt_LongitudYUnicidad(t *testing.T) {
	s1 := GenerateSalt()
	s2 := GenerateSalt()
	assert.Len(t, s1, 32) // 16 bytes en hex son 32 caracteres
	assert.NotEqual(t, s1, s2)
}
