package part

import (
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
)

type UpdateStockInput struct {
	ID       string
	Quantity int
}

type UpdateStockUseCase struct {
	partRepo repository.PartRepository
}

func NewUpdateStockUseCase(partRepo repository.PartRepository) *UpdateStockUseCase {
	return &UpdateStockUseCase{partRepo: partRepo}
}

func (u *UpdateStockUseCase) Execute(input UpdateStockInput) (*entity.Part, error) {
	part, err := u.partRepo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	part.AddToStock(input.Quantity)

	if err := u.partRepo.Update(part); err != nil {
		return nil, err
	}

	return part, nil
}
