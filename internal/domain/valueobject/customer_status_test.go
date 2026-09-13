package valueobject_test

import (
	"testing"

	"github.com/gabrielcamargo/oficina-api/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
)

func TestNewCustomerStatus_ShouldAcceptActive(t *testing.T) {
	status, err := valueobject.NewCustomerStatus("active")
	assert.NoError(t, err)
	assert.Equal(t, valueobject.CustomerStatusActive, status)
	assert.True(t, status.IsActive())
}

func TestNewCustomerStatus_ShouldAcceptInactive(t *testing.T) {
	status, err := valueobject.NewCustomerStatus("inactive")
	assert.NoError(t, err)
	assert.Equal(t, valueobject.CustomerStatusInactive, status)
	assert.False(t, status.IsActive())
}

func TestNewCustomerStatus_ShouldRejectUnknownValue(t *testing.T) {
	_, err := valueobject.NewCustomerStatus("blocked")
	assert.ErrorIs(t, err, valueobject.ErrInvalidCustomerStatus)
}
