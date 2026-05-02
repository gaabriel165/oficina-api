package valueobject_test

import (
	"testing"

	"github.com/gabrielcamargo/oficina-api/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
)

func TestNewCNPJ_ShouldReturnCNPJWhenValid(t *testing.T) {
	cnpj, err := valueobject.NewCNPJ("11.222.333/0001-81")
	assert.NoError(t, err)
	assert.Equal(t, "11222333000181", cnpj.Value())
}

func TestNewCNPJ_ShouldReturnCNPJWhenValidWithoutMask(t *testing.T) {
	cnpj, err := valueobject.NewCNPJ("11222333000181")
	assert.NoError(t, err)
	assert.Equal(t, "11222333000181", cnpj.Value())
}

func TestNewCNPJ_ShouldReturnErrorWhenAllDigitsSame(t *testing.T) {
	_, err := valueobject.NewCNPJ("11.111.111/1111-11")
	assert.ErrorIs(t, err, valueobject.ErrInvalidCNPJ)
}

func TestNewCNPJ_ShouldReturnErrorWhenTooShort(t *testing.T) {
	_, err := valueobject.NewCNPJ("11.222.333/0001")
	assert.ErrorIs(t, err, valueobject.ErrInvalidCNPJ)
}

func TestNewCNPJ_ShouldReturnErrorWhenInvalidCheckDigits(t *testing.T) {
	_, err := valueobject.NewCNPJ("11.222.333/0001-00")
	assert.ErrorIs(t, err, valueobject.ErrInvalidCNPJ)
}

func TestNewCNPJ_ShouldReturnErrorWhenEmpty(t *testing.T) {
	_, err := valueobject.NewCNPJ("")
	assert.ErrorIs(t, err, valueobject.ErrInvalidCNPJ)
}
