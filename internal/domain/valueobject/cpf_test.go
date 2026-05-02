package valueobject_test

import (
	"testing"

	"github.com/gabrielcamargo/oficina-api/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
)

func TestNewCPF_ShouldReturnCPFWhenValid(t *testing.T) {
	cpf, err := valueobject.NewCPF("529.982.247-25")
	assert.NoError(t, err)
	assert.Equal(t, "52998224725", cpf.Value())
}

func TestNewCPF_ShouldReturnCPFWhenValidWithoutMask(t *testing.T) {
	cpf, err := valueobject.NewCPF("52998224725")
	assert.NoError(t, err)
	assert.Equal(t, "52998224725", cpf.Value())
}

func TestNewCPF_ShouldReturnErrorWhenAllDigitsSame(t *testing.T) {
	_, err := valueobject.NewCPF("111.111.111-11")
	assert.ErrorIs(t, err, valueobject.ErrInvalidCPF)
}

func TestNewCPF_ShouldReturnErrorWhenTooShort(t *testing.T) {
	_, err := valueobject.NewCPF("123.456.789")
	assert.ErrorIs(t, err, valueobject.ErrInvalidCPF)
}

func TestNewCPF_ShouldReturnErrorWhenInvalidCheckDigits(t *testing.T) {
	_, err := valueobject.NewCPF("529.982.247-00")
	assert.ErrorIs(t, err, valueobject.ErrInvalidCPF)
}

func TestNewCPF_ShouldReturnErrorWhenEmpty(t *testing.T) {
	_, err := valueobject.NewCPF("")
	assert.ErrorIs(t, err, valueobject.ErrInvalidCPF)
}
