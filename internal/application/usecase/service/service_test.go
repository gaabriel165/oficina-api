package service_test

import (
	"testing"

	"github.com/gabrielcamargo/oficina-api/internal/application/mocks"
	"github.com/gabrielcamargo/oficina-api/internal/application/usecase/service"
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func makeService(t *testing.T) *entity.Service {
	t.Helper()
	s, _ := entity.NewService("Oil Change", "Full oil change", 150.00, 60)
	return s
}

func TestCreateService_ShouldCreateSuccessfully(t *testing.T) {
	serviceRepo := mocks.NewMockServiceRepository(t)

	serviceRepo.On("Create", mock.Anything).Return(nil)

	uc := service.NewCreateServiceUseCase(serviceRepo)
	result, err := uc.Execute(service.CreateServiceInput{
		Name:             "Oil Change",
		Description:      "Full oil change",
		LaborPrice:       150.00,
		EstimatedMinutes: 60,
	})

	assert.NoError(t, err)
	assert.Equal(t, "Oil Change", result.Name())
	assert.Equal(t, 150.00, result.LaborPrice())
	serviceRepo.AssertExpectations(t)
}

func TestCreateService_ShouldReturnErrorWhenPriceIsZero(t *testing.T) {
	serviceRepo := mocks.NewMockServiceRepository(t)

	uc := service.NewCreateServiceUseCase(serviceRepo)
	_, err := uc.Execute(service.CreateServiceInput{
		Name:             "Oil Change",
		LaborPrice:       0,
		EstimatedMinutes: 60,
	})

	assert.ErrorIs(t, err, entity.ErrServiceInvalidLaborPrice)
}

func TestCreateService_ShouldReturnErrorWhenEstimatedMinutesIsZero(t *testing.T) {
	serviceRepo := mocks.NewMockServiceRepository(t)

	uc := service.NewCreateServiceUseCase(serviceRepo)
	_, err := uc.Execute(service.CreateServiceInput{
		Name:             "Oil Change",
		LaborPrice:       150.00,
		EstimatedMinutes: 0,
	})

	assert.ErrorIs(t, err, entity.ErrServiceInvalidEstimatedMin)
}

func TestUpdateService_ShouldUpdateSuccessfully(t *testing.T) {
	serviceRepo := mocks.NewMockServiceRepository(t)

	serviceRepo.On("FindByID", "service-id").Return(makeService(t), nil)
	serviceRepo.On("Update", mock.Anything).Return(nil)

	uc := service.NewUpdateServiceUseCase(serviceRepo)
	result, err := uc.Execute(service.UpdateServiceInput{
		ID:               "service-id",
		Name:             "Premium Oil Change",
		Description:      "Premium full oil change",
		LaborPrice:       200.00,
		EstimatedMinutes: 90,
	})

	assert.NoError(t, err)
	assert.Equal(t, "Premium Oil Change", result.Name())
	assert.Equal(t, 200.00, result.LaborPrice())
}

func TestUpdateService_ShouldReturnErrorWhenNotFound(t *testing.T) {
	serviceRepo := mocks.NewMockServiceRepository(t)

	serviceRepo.On("FindByID", "service-id").Return(nil, repository.ErrServiceNotFound)

	uc := service.NewUpdateServiceUseCase(serviceRepo)
	_, err := uc.Execute(service.UpdateServiceInput{
		ID:               "service-id",
		Name:             "Oil Change",
		LaborPrice:       150.00,
		EstimatedMinutes: 60,
	})

	assert.ErrorIs(t, err, repository.ErrServiceNotFound)
}

func TestDeleteService_ShouldDeleteSuccessfully(t *testing.T) {
	serviceRepo := mocks.NewMockServiceRepository(t)

	serviceRepo.On("FindByID", "service-id").Return(makeService(t), nil)
	serviceRepo.On("Delete", "service-id").Return(nil)

	uc := service.NewDeleteServiceUseCase(serviceRepo)
	err := uc.Execute("service-id")

	assert.NoError(t, err)
	serviceRepo.AssertExpectations(t)
}

func TestDeleteService_ShouldReturnErrorWhenNotFound(t *testing.T) {
	serviceRepo := mocks.NewMockServiceRepository(t)

	serviceRepo.On("FindByID", "service-id").Return(nil, repository.ErrServiceNotFound)

	uc := service.NewDeleteServiceUseCase(serviceRepo)
	err := uc.Execute("service-id")

	assert.ErrorIs(t, err, repository.ErrServiceNotFound)
}

func TestGetService_ShouldReturnService(t *testing.T) {
	serviceRepo := mocks.NewMockServiceRepository(t)

	serviceRepo.On("FindByID", "service-id").Return(makeService(t), nil)

	uc := service.NewGetServiceUseCase(serviceRepo)
	result, err := uc.Execute("service-id")

	assert.NoError(t, err)
	assert.Equal(t, "Oil Change", result.Name())
}

func TestListServices_ShouldReturnAll(t *testing.T) {
	serviceRepo := mocks.NewMockServiceRepository(t)

	serviceRepo.On("FindAll").Return([]*entity.Service{makeService(t), makeService(t)}, nil)

	uc := service.NewListServicesUseCase(serviceRepo)
	results, err := uc.Execute()

	assert.NoError(t, err)
	assert.Len(t, results, 2)
}
