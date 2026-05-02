package entity

import "github.com/google/uuid"

type ServiceOrderPart struct {
	id             string
	serviceOrderID string
	partID         string
	partName       string
	quantity       int
	unitPrice      float64
}

func newServiceOrderPart(serviceOrderID string, part *Part, quantity int) *ServiceOrderPart {
	return &ServiceOrderPart{
		id:             uuid.NewString(),
		serviceOrderID: serviceOrderID,
		partID:         part.ID(),
		partName:       part.Name(),
		quantity:       quantity,
		unitPrice:      part.UnitPrice(),
	}
}

func RestoreServiceOrderPart(id, serviceOrderID, partID, partName string, quantity int, unitPrice float64) *ServiceOrderPart {
	return &ServiceOrderPart{
		id:             id,
		serviceOrderID: serviceOrderID,
		partID:         partID,
		partName:       partName,
		quantity:       quantity,
		unitPrice:      unitPrice,
	}
}

func (p *ServiceOrderPart) ID() string             { return p.id }
func (p *ServiceOrderPart) ServiceOrderID() string { return p.serviceOrderID }
func (p *ServiceOrderPart) PartID() string         { return p.partID }
func (p *ServiceOrderPart) PartName() string       { return p.partName }
func (p *ServiceOrderPart) Quantity() int          { return p.quantity }
func (p *ServiceOrderPart) UnitPrice() float64     { return p.unitPrice }
func (p *ServiceOrderPart) TotalPrice() float64    { return p.unitPrice * float64(p.quantity) }
