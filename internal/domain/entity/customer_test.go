package entity_test

import (
	"testing"

	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
)

func TestNewCustomer_ShouldCreateWithCPF(t *testing.T) {
	customer, err := entity.NewCustomer("John Doe", "529.982.247-25", "11999999999", "john@email.com")
	assert.NoError(t, err)
	assert.NotEmpty(t, customer.ID())
	assert.Equal(t, "John Doe", customer.Name())
	assert.Equal(t, "52998224725", customer.Document())
	assert.Equal(t, entity.DocumentTypeCPF, customer.DocumentType())
}

func TestNewCustomer_ShouldCreateWithCNPJ(t *testing.T) {
	customer, err := entity.NewCustomer("Acme Corp", "11.222.333/0001-81", "1133334444", "acme@email.com")
	assert.NoError(t, err)
	assert.Equal(t, entity.DocumentTypeCNPJ, customer.DocumentType())
	assert.Equal(t, "11222333000181", customer.Document())
}

func TestNewCustomer_ShouldReturnErrorWhenNameEmpty(t *testing.T) {
	_, err := entity.NewCustomer("", "529.982.247-25", "11999999999", "john@email.com")
	assert.ErrorIs(t, err, entity.ErrCustomerNameRequired)
}

func TestNewCustomer_ShouldReturnErrorWhenPhoneEmpty(t *testing.T) {
	_, err := entity.NewCustomer("John Doe", "529.982.247-25", "", "john@email.com")
	assert.ErrorIs(t, err, entity.ErrCustomerPhoneRequired)
}

func TestNewCustomer_ShouldReturnErrorWhenEmailEmpty(t *testing.T) {
	_, err := entity.NewCustomer("John Doe", "529.982.247-25", "11999999999", "")
	assert.ErrorIs(t, err, entity.ErrCustomerEmailRequired)
}

func TestNewCustomer_ShouldReturnErrorWhenInvalidDocument(t *testing.T) {
	_, err := entity.NewCustomer("John Doe", "000.000.000-00", "11999999999", "john@email.com")
	assert.ErrorIs(t, err, entity.ErrInvalidDocument)
}

func TestCustomer_Update_ShouldUpdateFields(t *testing.T) {
	customer, _ := entity.NewCustomer("John Doe", "529.982.247-25", "11999999999", "john@email.com")
	err := customer.Update("Jane Doe", "11988888888", "jane@email.com")
	assert.NoError(t, err)
	assert.Equal(t, "Jane Doe", customer.Name())
	assert.Equal(t, "11988888888", customer.Phone())
	assert.Equal(t, "jane@email.com", customer.Email())
}

func TestCustomer_Update_ShouldReturnErrorWhenNameEmpty(t *testing.T) {
	customer, _ := entity.NewCustomer("John Doe", "529.982.247-25", "11999999999", "john@email.com")
	err := customer.Update("", "11988888888", "jane@email.com")
	assert.ErrorIs(t, err, entity.ErrCustomerNameRequired)
}

func TestCustomer_Update_ShouldReturnErrorWhenPhoneEmpty(t *testing.T) {
	customer, _ := entity.NewCustomer("John Doe", "529.982.247-25", "11999999999", "john@email.com")
	err := customer.Update("Jane Doe", "", "jane@email.com")
	assert.ErrorIs(t, err, entity.ErrCustomerPhoneRequired)
}

func TestCustomer_Update_ShouldReturnErrorWhenEmailEmpty(t *testing.T) {
	customer, _ := entity.NewCustomer("John Doe", "529.982.247-25", "11999999999", "john@email.com")
	err := customer.Update("Jane Doe", "11988888888", "")
	assert.ErrorIs(t, err, entity.ErrCustomerEmailRequired)
}

func TestNewCustomer_ShouldStartActive(t *testing.T) {
	customer, err := entity.NewCustomer("John Doe", "529.982.247-25", "11999999999", "john@email.com")
	assert.NoError(t, err)
	assert.True(t, customer.IsActive())
}

func TestCustomerDeactivate_ShouldChangeStatusToInactive(t *testing.T) {
	customer, _ := entity.NewCustomer("John Doe", "529.982.247-25", "11999999999", "john@email.com")
	customer.Deactivate()
	assert.False(t, customer.IsActive())
	assert.Equal(t, valueobject.CustomerStatusInactive, customer.Status())
}

func TestCustomerActivate_ShouldChangeStatusToActive(t *testing.T) {
	customer, _ := entity.NewCustomer("John Doe", "529.982.247-25", "11999999999", "john@email.com")
	customer.Deactivate()
	customer.Activate()
	assert.True(t, customer.IsActive())
}
