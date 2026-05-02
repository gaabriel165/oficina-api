package entity_test

import (
	"testing"
	"time"

	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewVehicle_ShouldCreateWithValidData(t *testing.T) {
	vehicle, err := entity.NewVehicle("customer-id", "ABC-1234", "Toyota", "Corolla", 2020)
	assert.NoError(t, err)
	assert.NotEmpty(t, vehicle.ID())
	assert.Equal(t, "ABC1234", vehicle.Plate())
	assert.Equal(t, "Toyota", vehicle.Brand())
	assert.Equal(t, "Corolla", vehicle.Model())
	assert.Equal(t, 2020, vehicle.Year())
}

func TestNewVehicle_ShouldCreateWithMercosulPlate(t *testing.T) {
	vehicle, err := entity.NewVehicle("customer-id", "ABC1D23", "Honda", "Civic", 2022)
	assert.NoError(t, err)
	assert.Equal(t, "ABC1D23", vehicle.Plate())
}

func TestNewVehicle_ShouldReturnErrorWhenCustomerIDEmpty(t *testing.T) {
	_, err := entity.NewVehicle("", "ABC-1234", "Toyota", "Corolla", 2020)
	assert.ErrorIs(t, err, entity.ErrVehicleCustomerRequired)
}

func TestNewVehicle_ShouldReturnErrorWhenBrandEmpty(t *testing.T) {
	_, err := entity.NewVehicle("customer-id", "ABC-1234", "", "Corolla", 2020)
	assert.ErrorIs(t, err, entity.ErrVehicleBrandRequired)
}

func TestNewVehicle_ShouldReturnErrorWhenInvalidYear(t *testing.T) {
	_, err := entity.NewVehicle("customer-id", "ABC-1234", "Toyota", "Corolla", 1800)
	assert.ErrorIs(t, err, entity.ErrVehicleInvalidYear)
}

func TestNewVehicle_ShouldReturnErrorWhenFutureYear(t *testing.T) {
	_, err := entity.NewVehicle("customer-id", "ABC-1234", "Toyota", "Corolla", time.Now().Year()+5)
	assert.ErrorIs(t, err, entity.ErrVehicleInvalidYear)
}

func TestNewVehicle_ShouldReturnErrorWhenInvalidPlate(t *testing.T) {
	_, err := entity.NewVehicle("customer-id", "INVALID", "Toyota", "Corolla", 2020)
	assert.Error(t, err)
}

func TestVehicle_Update_ShouldUpdateFields(t *testing.T) {
	vehicle, _ := entity.NewVehicle("customer-id", "ABC-1234", "Toyota", "Corolla", 2020)
	err := vehicle.Update("Honda", "Civic", 2022)
	assert.NoError(t, err)
	assert.Equal(t, "Honda", vehicle.Brand())
	assert.Equal(t, "Civic", vehicle.Model())
	assert.Equal(t, 2022, vehicle.Year())
}

func TestVehicle_Update_ShouldReturnErrorWhenBrandEmpty(t *testing.T) {
	vehicle, _ := entity.NewVehicle("customer-id", "ABC-1234", "Toyota", "Corolla", 2020)
	err := vehicle.Update("", "Civic", 2022)
	assert.ErrorIs(t, err, entity.ErrVehicleBrandRequired)
}

func TestVehicle_Update_ShouldReturnErrorWhenModelEmpty(t *testing.T) {
	vehicle, _ := entity.NewVehicle("customer-id", "ABC-1234", "Toyota", "Corolla", 2020)
	err := vehicle.Update("Honda", "", 2022)
	assert.ErrorIs(t, err, entity.ErrVehicleModelRequired)
}

func TestVehicle_Update_ShouldReturnErrorWhenInvalidYear(t *testing.T) {
	vehicle, _ := entity.NewVehicle("customer-id", "ABC-1234", "Toyota", "Corolla", 2020)
	err := vehicle.Update("Honda", "Civic", 1800)
	assert.ErrorIs(t, err, entity.ErrVehicleInvalidYear)
}
