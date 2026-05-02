package entity

import "github.com/google/uuid"

type ServiceOrderItem struct {
	id             string
	serviceOrderID string
	serviceID      string
	serviceName    string
	laborPrice     float64
}

func newServiceOrderItem(serviceOrderID string, service *Service) *ServiceOrderItem {
	return &ServiceOrderItem{
		id:             uuid.NewString(),
		serviceOrderID: serviceOrderID,
		serviceID:      service.ID(),
		serviceName:    service.Name(),
		laborPrice:     service.LaborPrice(),
	}
}

func RestoreServiceOrderItem(id, serviceOrderID, serviceID, serviceName string, laborPrice float64) *ServiceOrderItem {
	return &ServiceOrderItem{
		id:             id,
		serviceOrderID: serviceOrderID,
		serviceID:      serviceID,
		serviceName:    serviceName,
		laborPrice:     laborPrice,
	}
}

func (i *ServiceOrderItem) ID() string             { return i.id }
func (i *ServiceOrderItem) ServiceOrderID() string { return i.serviceOrderID }
func (i *ServiceOrderItem) ServiceID() string      { return i.serviceID }
func (i *ServiceOrderItem) ServiceName() string    { return i.serviceName }
func (i *ServiceOrderItem) LaborPrice() float64    { return i.laborPrice }
