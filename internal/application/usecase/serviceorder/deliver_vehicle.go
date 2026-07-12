package serviceorder

import (
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
)

type DeliverVehicleUseCase struct {
	orderRepo repository.ServiceOrderRepository
	notifier  StatusNotifier
}

func NewDeliverVehicleUseCase(orderRepo repository.ServiceOrderRepository, notifier StatusNotifier) *DeliverVehicleUseCase {
	return &DeliverVehicleUseCase{orderRepo: orderRepo, notifier: notifier}
}

func (u *DeliverVehicleUseCase) Execute(orderID string) (*entity.ServiceOrder, error) {
	order, err := u.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, err
	}

	if err := order.Deliver(); err != nil {
		return nil, err
	}

	if err := u.orderRepo.Update(order); err != nil {
		return nil, err
	}

	u.notifier.Notify(order)

	return order, nil
}
