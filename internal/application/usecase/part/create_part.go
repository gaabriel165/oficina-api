package part

import (
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
)

type CreatePartInput struct {
	Name          string
	Description   string
	UnitPrice     float64
	StockQuantity int
}

type CreatePartUseCase struct {
	partRepo repository.PartRepository
}

func NewCreatePartUseCase(partRepo repository.PartRepository) *CreatePartUseCase {
	return &CreatePartUseCase{partRepo: partRepo}
}

func (u *CreatePartUseCase) Execute(input CreatePartInput) (*entity.Part, error) {
	part, err := entity.NewPart(input.Name, input.Description, input.UnitPrice, input.StockQuantity)
	if err != nil {
		return nil, err
	}

	if err := u.partRepo.Create(part); err != nil {
		return nil, err
	}

	return part, nil
}
