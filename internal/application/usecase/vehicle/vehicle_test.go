package vehicle_test

import (
	"testing"

	"github.com/gabrielcamargo/oficina-api/internal/application/mocks"
	"github.com/gabrielcamargo/oficina-api/internal/application/usecase"
	"github.com/gabrielcamargo/oficina-api/internal/application/usecase/vehicle"
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func makeCustomer(t *testing.T) *entity.Customer {
	t.Helper()
	c, _ := entity.NewCustomer("John Doe", "529.982.247-25", "11999999999", "john@email.com")
	return c
}

func makeVehicle(t *testing.T) *entity.Vehicle {
	t.Helper()
	v, _ := entity.NewVehicle("customer-id", "ABC-1234", "Toyota", "Corolla", 2020)
	return v
}

func TestCreateVehicle_ShouldCreateSuccessfully(t *testing.T) {
	vehicleRepo := mocks.NewMockVehicleRepository(t)
	customerRepo := mocks.NewMockCustomerRepository(t)

	customerRepo.On("FindByID", "customer-id").Return(makeCustomer(t), nil)
	vehicleRepo.On("FindByPlate", "ABC1234").Return(nil, repository.ErrVehicleNotFound)
	vehicleRepo.On("Create", mock.Anything).Return(nil)

	uc := vehicle.NewCreateVehicleUseCase(vehicleRepo, customerRepo)
	result, err := uc.Execute(vehicle.CreateVehicleInput{
		CustomerID: "customer-id",
		Plate:      "ABC-1234",
		Brand:      "Toyota",
		Model:      "Corolla",
		Year:       2020,
	})

	assert.NoError(t, err)
	assert.Equal(t, "Toyota", result.Brand())
	assert.Equal(t, "ABC1234", result.Plate())
	vehicleRepo.AssertExpectations(t)
}

func TestCreateVehicle_ShouldReturnErrorWhenCustomerNotFound(t *testing.T) {
	vehicleRepo := mocks.NewMockVehicleRepository(t)
	customerRepo := mocks.NewMockCustomerRepository(t)

	customerRepo.On("FindByID", "customer-id").Return(nil, repository.ErrCustomerNotFound)

	uc := vehicle.NewCreateVehicleUseCase(vehicleRepo, customerRepo)
	_, err := uc.Execute(vehicle.CreateVehicleInput{
		CustomerID: "customer-id",
		Plate:      "ABC-1234",
		Brand:      "Toyota",
		Model:      "Corolla",
		Year:       2020,
	})

	assert.ErrorIs(t, err, repository.ErrCustomerNotFound)
}

func TestCreateVehicle_ShouldReturnErrorWhenPlateAlreadyRegistered(t *testing.T) {
	vehicleRepo := mocks.NewMockVehicleRepository(t)
	customerRepo := mocks.NewMockCustomerRepository(t)

	customerRepo.On("FindByID", "customer-id").Return(makeCustomer(t), nil)
	vehicleRepo.On("FindByPlate", "ABC1234").Return(makeVehicle(t), nil)

	uc := vehicle.NewCreateVehicleUseCase(vehicleRepo, customerRepo)
	_, err := uc.Execute(vehicle.CreateVehicleInput{
		CustomerID: "customer-id",
		Plate:      "ABC-1234",
		Brand:      "Toyota",
		Model:      "Corolla",
		Year:       2020,
	})

	assert.ErrorIs(t, err, usecase.ErrPlateAlreadyRegistered)
}

func TestUpdateVehicle_ShouldUpdateSuccessfully(t *testing.T) {
	vehicleRepo := mocks.NewMockVehicleRepository(t)

	vehicleRepo.On("FindByID", "vehicle-id").Return(makeVehicle(t), nil)
	vehicleRepo.On("Update", mock.Anything).Return(nil)

	uc := vehicle.NewUpdateVehicleUseCase(vehicleRepo)
	result, err := uc.Execute(vehicle.UpdateVehicleInput{
		ID:    "vehicle-id",
		Brand: "Honda",
		Model: "Civic",
		Year:  2022,
	})

	assert.NoError(t, err)
	assert.Equal(t, "Honda", result.Brand())
	assert.Equal(t, 2022, result.Year())
}

func TestUpdateVehicle_ShouldReturnErrorWhenNotFound(t *testing.T) {
	vehicleRepo := mocks.NewMockVehicleRepository(t)

	vehicleRepo.On("FindByID", "vehicle-id").Return(nil, repository.ErrVehicleNotFound)

	uc := vehicle.NewUpdateVehicleUseCase(vehicleRepo)
	_, err := uc.Execute(vehicle.UpdateVehicleInput{
		ID:    "vehicle-id",
		Brand: "Honda",
		Model: "Civic",
		Year:  2022,
	})

	assert.ErrorIs(t, err, repository.ErrVehicleNotFound)
}

func TestDeleteVehicle_ShouldDeleteSuccessfully(t *testing.T) {
	vehicleRepo := mocks.NewMockVehicleRepository(t)

	vehicleRepo.On("FindByID", "vehicle-id").Return(makeVehicle(t), nil)
	vehicleRepo.On("Delete", "vehicle-id").Return(nil)

	uc := vehicle.NewDeleteVehicleUseCase(vehicleRepo)
	err := uc.Execute("vehicle-id")

	assert.NoError(t, err)
	vehicleRepo.AssertExpectations(t)
}

func TestDeleteVehicle_ShouldReturnErrorWhenNotFound(t *testing.T) {
	vehicleRepo := mocks.NewMockVehicleRepository(t)

	vehicleRepo.On("FindByID", "vehicle-id").Return(nil, repository.ErrVehicleNotFound)

	uc := vehicle.NewDeleteVehicleUseCase(vehicleRepo)
	err := uc.Execute("vehicle-id")

	assert.ErrorIs(t, err, repository.ErrVehicleNotFound)
}

func TestGetVehicle_ShouldReturnVehicle(t *testing.T) {
	vehicleRepo := mocks.NewMockVehicleRepository(t)

	vehicleRepo.On("FindByID", "vehicle-id").Return(makeVehicle(t), nil)

	uc := vehicle.NewGetVehicleUseCase(vehicleRepo)
	result, err := uc.Execute("vehicle-id")

	assert.NoError(t, err)
	assert.Equal(t, "Toyota", result.Brand())
}

func TestListVehicles_ShouldReturnAll(t *testing.T) {
	vehicleRepo := mocks.NewMockVehicleRepository(t)

	vehicleRepo.On("FindAll").Return([]*entity.Vehicle{makeVehicle(t), makeVehicle(t)}, nil)

	uc := vehicle.NewListVehiclesUseCase(vehicleRepo)
	results, err := uc.Execute()

	assert.NoError(t, err)
	assert.Len(t, results, 2)
}
