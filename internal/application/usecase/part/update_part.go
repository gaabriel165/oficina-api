package part

import (
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
)

type UpdatePartInput struct {
	ID          string
	Name        string
	Description string
	UnitPrice   float64
}

type UpdatePartUseCase struct {
	partRepo repository.PartRepository
}

func NewUpdatePartUseCase(partRepo repository.PartRepository) *UpdatePartUseCase {
	return &UpdatePartUseCase{partRepo: partRepo}
}

func (u *UpdatePartUseCase) Execute(input UpdatePartInput) (*entity.Part, error) {
	part, err := u.partRepo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	if err := part.Update(input.Name, input.Description, input.UnitPrice); err != nil {
		return nil, err
	}

	if err := u.partRepo.Update(part); err != nil {
		return nil, err
	}

	return part, nil
}
