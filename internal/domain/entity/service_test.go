package entity_test

import (
	"testing"

	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewService_ShouldCreateWithValidData(t *testing.T) {
	svc, err := entity.NewService("Oil Change", "Full oil change", 150.00, 60)
	assert.NoError(t, err)
	assert.NotEmpty(t, svc.ID())
	assert.Equal(t, "Oil Change", svc.Name())
	assert.Equal(t, 150.00, svc.LaborPrice())
	assert.Equal(t, 60, svc.EstimatedMinutes())
}

func TestNewService_ShouldReturnErrorWhenNameEmpty(t *testing.T) {
	_, err := entity.NewService("", "desc", 150.00, 60)
	assert.ErrorIs(t, err, entity.ErrServiceNameRequired)
}

func TestNewService_ShouldReturnErrorWhenPriceZero(t *testing.T) {
	_, err := entity.NewService("Oil Change", "desc", 0, 60)
	assert.ErrorIs(t, err, entity.ErrServiceInvalidLaborPrice)
}

func TestNewService_ShouldReturnErrorWhenEstimatedMinutesZero(t *testing.T) {
	_, err := entity.NewService("Oil Change", "desc", 150.00, 0)
	assert.ErrorIs(t, err, entity.ErrServiceInvalidEstimatedMin)
}

func TestService_Update_ShouldUpdateFields(t *testing.T) {
	svc, _ := entity.NewService("Oil Change", "desc", 150.00, 60)
	err := svc.Update("Premium Oil Change", "premium desc", 200.00, 90)
	assert.NoError(t, err)
	assert.Equal(t, "Premium Oil Change", svc.Name())
	assert.Equal(t, 200.00, svc.LaborPrice())
	assert.Equal(t, 90, svc.EstimatedMinutes())
}

func TestService_Update_ShouldReturnErrorWhenNameEmpty(t *testing.T) {
	svc, _ := entity.NewService("Oil Change", "desc", 150.00, 60)
	err := svc.Update("", "desc", 150.00, 60)
	assert.ErrorIs(t, err, entity.ErrServiceNameRequired)
}

func TestService_Update_ShouldReturnErrorWhenPriceZero(t *testing.T) {
	svc, _ := entity.NewService("Oil Change", "desc", 150.00, 60)
	err := svc.Update("Oil Change", "desc", 0, 60)
	assert.ErrorIs(t, err, entity.ErrServiceInvalidLaborPrice)
}

func TestService_Update_ShouldReturnErrorWhenEstimatedMinutesZero(t *testing.T) {
	svc, _ := entity.NewService("Oil Change", "desc", 150.00, 60)
	err := svc.Update("Oil Change", "desc", 150.00, 0)
	assert.ErrorIs(t, err, entity.ErrServiceInvalidEstimatedMin)
}
