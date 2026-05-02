package entity_test

import (
	"testing"

	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewPart_ShouldCreateWithValidData(t *testing.T) {
	part, err := entity.NewPart("Oil Filter", "Standard oil filter", 29.90, 10)
	assert.NoError(t, err)
	assert.NotEmpty(t, part.ID())
	assert.Equal(t, "Oil Filter", part.Name())
	assert.Equal(t, 29.90, part.UnitPrice())
	assert.Equal(t, 10, part.StockQuantity())
}

func TestNewPart_ShouldReturnErrorWhenNameEmpty(t *testing.T) {
	_, err := entity.NewPart("", "desc", 29.90, 10)
	assert.ErrorIs(t, err, entity.ErrPartNameRequired)
}

func TestNewPart_ShouldReturnErrorWhenPriceZero(t *testing.T) {
	_, err := entity.NewPart("Oil Filter", "desc", 0, 10)
	assert.ErrorIs(t, err, entity.ErrPartInvalidPrice)
}

func TestNewPart_ShouldReturnErrorWhenNegativeStock(t *testing.T) {
	_, err := entity.NewPart("Oil Filter", "desc", 29.90, -1)
	assert.ErrorIs(t, err, entity.ErrPartInvalidStock)
}

func TestPart_DebitStock_ShouldDecrementQuantity(t *testing.T) {
	part, _ := entity.NewPart("Oil Filter", "desc", 29.90, 10)
	err := part.DebitStock(3)
	assert.NoError(t, err)
	assert.Equal(t, 7, part.StockQuantity())
}

func TestPart_DebitStock_ShouldReturnErrorWhenInsufficientStock(t *testing.T) {
	part, _ := entity.NewPart("Oil Filter", "desc", 29.90, 2)
	err := part.DebitStock(5)
	assert.ErrorIs(t, err, entity.ErrPartInsufficientStock)
}

func TestPart_AddToStock_ShouldIncrementQuantity(t *testing.T) {
	part, _ := entity.NewPart("Oil Filter", "desc", 29.90, 5)
	part.AddToStock(3)
	assert.Equal(t, 8, part.StockQuantity())
}

func TestPart_Update_ShouldUpdateFields(t *testing.T) {
	part, _ := entity.NewPart("Oil Filter", "desc", 29.90, 10)
	err := part.Update("Premium Filter", "premium desc", 49.90)
	assert.NoError(t, err)
	assert.Equal(t, "Premium Filter", part.Name())
	assert.Equal(t, 49.90, part.UnitPrice())
}

func TestPart_Update_ShouldReturnErrorWhenNameEmpty(t *testing.T) {
	part, _ := entity.NewPart("Oil Filter", "desc", 29.90, 10)
	err := part.Update("", "desc", 29.90)
	assert.ErrorIs(t, err, entity.ErrPartNameRequired)
}

func TestPart_Update_ShouldReturnErrorWhenPriceZero(t *testing.T) {
	part, _ := entity.NewPart("Oil Filter", "desc", 29.90, 10)
	err := part.Update("Oil Filter", "desc", 0)
	assert.ErrorIs(t, err, entity.ErrPartInvalidPrice)
}
