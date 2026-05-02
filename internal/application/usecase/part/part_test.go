package part_test

import (
	"testing"

	"github.com/gabrielcamargo/oficina-api/internal/application/mocks"
	"github.com/gabrielcamargo/oficina-api/internal/application/usecase/part"
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func makePart(t *testing.T) *entity.Part {
	t.Helper()
	p, _ := entity.NewPart("Oil Filter", "Standard filter", 29.90, 50)
	return p
}

func TestCreatePart_ShouldCreateSuccessfully(t *testing.T) {
	partRepo := mocks.NewMockPartRepository(t)

	partRepo.On("Create", mock.Anything).Return(nil)

	uc := part.NewCreatePartUseCase(partRepo)
	result, err := uc.Execute(part.CreatePartInput{
		Name:          "Oil Filter",
		Description:   "Standard filter",
		UnitPrice:     29.90,
		StockQuantity: 50,
	})

	assert.NoError(t, err)
	assert.Equal(t, "Oil Filter", result.Name())
	assert.Equal(t, 50, result.StockQuantity())
	partRepo.AssertExpectations(t)
}

func TestCreatePart_ShouldReturnErrorWhenPriceIsZero(t *testing.T) {
	partRepo := mocks.NewMockPartRepository(t)

	uc := part.NewCreatePartUseCase(partRepo)
	_, err := uc.Execute(part.CreatePartInput{
		Name:      "Oil Filter",
		UnitPrice: 0,
	})

	assert.ErrorIs(t, err, entity.ErrPartInvalidPrice)
}

func TestUpdatePart_ShouldUpdateSuccessfully(t *testing.T) {
	partRepo := mocks.NewMockPartRepository(t)

	partRepo.On("FindByID", "part-id").Return(makePart(t), nil)
	partRepo.On("Update", mock.Anything).Return(nil)

	uc := part.NewUpdatePartUseCase(partRepo)
	result, err := uc.Execute(part.UpdatePartInput{
		ID:          "part-id",
		Name:        "Premium Oil Filter",
		Description: "Premium filter",
		UnitPrice:   39.90,
	})

	assert.NoError(t, err)
	assert.Equal(t, "Premium Oil Filter", result.Name())
	assert.Equal(t, 39.90, result.UnitPrice())
}

func TestUpdatePart_ShouldReturnErrorWhenNotFound(t *testing.T) {
	partRepo := mocks.NewMockPartRepository(t)

	partRepo.On("FindByID", "part-id").Return(nil, repository.ErrPartNotFound)

	uc := part.NewUpdatePartUseCase(partRepo)
	_, err := uc.Execute(part.UpdatePartInput{
		ID:        "part-id",
		Name:      "Filter",
		UnitPrice: 29.90,
	})

	assert.ErrorIs(t, err, repository.ErrPartNotFound)
}

func TestDeletePart_ShouldDeleteSuccessfully(t *testing.T) {
	partRepo := mocks.NewMockPartRepository(t)

	partRepo.On("FindByID", "part-id").Return(makePart(t), nil)
	partRepo.On("Delete", "part-id").Return(nil)

	uc := part.NewDeletePartUseCase(partRepo)
	err := uc.Execute("part-id")

	assert.NoError(t, err)
	partRepo.AssertExpectations(t)
}

func TestDeletePart_ShouldReturnErrorWhenNotFound(t *testing.T) {
	partRepo := mocks.NewMockPartRepository(t)

	partRepo.On("FindByID", "part-id").Return(nil, repository.ErrPartNotFound)

	uc := part.NewDeletePartUseCase(partRepo)
	err := uc.Execute("part-id")

	assert.ErrorIs(t, err, repository.ErrPartNotFound)
}

func TestGetPart_ShouldReturnPart(t *testing.T) {
	partRepo := mocks.NewMockPartRepository(t)

	partRepo.On("FindByID", "part-id").Return(makePart(t), nil)

	uc := part.NewGetPartUseCase(partRepo)
	result, err := uc.Execute("part-id")

	assert.NoError(t, err)
	assert.Equal(t, "Oil Filter", result.Name())
}

func TestListParts_ShouldReturnAll(t *testing.T) {
	partRepo := mocks.NewMockPartRepository(t)

	partRepo.On("FindAll").Return([]*entity.Part{makePart(t), makePart(t)}, nil)

	uc := part.NewListPartsUseCase(partRepo)
	results, err := uc.Execute()

	assert.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestUpdateStock_ShouldAddToStock(t *testing.T) {
	partRepo := mocks.NewMockPartRepository(t)
	p := makePart(t)

	partRepo.On("FindByID", "part-id").Return(p, nil)
	partRepo.On("Update", mock.Anything).Return(nil)

	uc := part.NewUpdateStockUseCase(partRepo)
	result, err := uc.Execute(part.UpdateStockInput{
		ID:       "part-id",
		Quantity: 20,
	})

	assert.NoError(t, err)
	assert.Equal(t, 70, result.StockQuantity())
}
