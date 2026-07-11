package serviceorder

import (
	"sort"

	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
	"github.com/gabrielcamargo/oficina-api/internal/domain/valueobject"
)

var activeListingStatuses = []valueobject.OrderStatus{
	valueobject.OrderStatusReceived,
	valueobject.OrderStatusInDiagnosis,
	valueobject.OrderStatusWaitingApproval,
	valueobject.OrderStatusInExecution,
}

type ListServiceOrdersUseCase struct {
	orderRepo repository.ServiceOrderRepository
}

func NewListServiceOrdersUseCase(orderRepo repository.ServiceOrderRepository) *ListServiceOrdersUseCase {
	return &ListServiceOrdersUseCase{orderRepo: orderRepo}
}

func (u *ListServiceOrdersUseCase) Execute() ([]*entity.ServiceOrder, error) {
	orders, err := u.orderRepo.FindByStatuses(activeListingStatuses)
	if err != nil {
		return nil, err
	}

	sort.SliceStable(orders, func(i, j int) bool {
		left, right := orders[i], orders[j]
		if left.Status().ListingPriority() != right.Status().ListingPriority() {
			return left.Status().ListingPriority() < right.Status().ListingPriority()
		}
		return left.CreatedAt().Before(right.CreatedAt())
	})

	return orders, nil
}
