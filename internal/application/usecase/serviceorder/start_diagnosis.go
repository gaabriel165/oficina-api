package serviceorder

import (
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
)

type StartDiagnosisUseCase struct {
	orderRepo repository.ServiceOrderRepository
}

func NewStartDiagnosisUseCase(orderRepo repository.ServiceOrderRepository) *StartDiagnosisUseCase {
	return &StartDiagnosisUseCase{orderRepo: orderRepo}
}

func (u *StartDiagnosisUseCase) Execute(orderID string) (*entity.ServiceOrder, error) {
	order, err := u.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, err
	}

	if err := order.StartDiagnosis(); err != nil {
		return nil, err
	}

	if err := u.orderRepo.Update(order); err != nil {
		return nil, err
	}

	return order, nil
}
