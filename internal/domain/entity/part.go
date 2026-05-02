package entity

import (
	"time"

	"github.com/google/uuid"
)

type Part struct {
	id            string
	name          string
	description   string
	unitPrice     float64
	stockQuantity int
	createdAt     time.Time
	updatedAt     time.Time
}

func NewPart(name, description string, unitPrice float64, stockQuantity int) (*Part, error) {
	if name == "" {
		return nil, ErrPartNameRequired
	}
	if unitPrice <= 0 {
		return nil, ErrPartInvalidPrice
	}
	if stockQuantity < 0 {
		return nil, ErrPartInvalidStock
	}

	now := time.Now()
	return &Part{
		id:            uuid.NewString(),
		name:          name,
		description:   description,
		unitPrice:     unitPrice,
		stockQuantity: stockQuantity,
		createdAt:     now,
		updatedAt:     now,
	}, nil
}

func RestorePart(id, name, description string, unitPrice float64, stockQuantity int, createdAt, updatedAt time.Time) *Part {
	return &Part{
		id:            id,
		name:          name,
		description:   description,
		unitPrice:     unitPrice,
		stockQuantity: stockQuantity,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
	}
}

func (p *Part) ID() string            { return p.id }
func (p *Part) Name() string          { return p.name }
func (p *Part) Description() string   { return p.description }
func (p *Part) UnitPrice() float64    { return p.unitPrice }
func (p *Part) StockQuantity() int    { return p.stockQuantity }
func (p *Part) CreatedAt() time.Time  { return p.createdAt }
func (p *Part) UpdatedAt() time.Time  { return p.updatedAt }

func (p *Part) HasSufficientStock(quantity int) bool {
	return p.stockQuantity >= quantity
}

func (p *Part) DebitStock(quantity int) error {
	if !p.HasSufficientStock(quantity) {
		return ErrPartInsufficientStock
	}
	p.stockQuantity -= quantity
	p.updatedAt = time.Now()
	return nil
}

func (p *Part) AddToStock(quantity int) {
	p.stockQuantity += quantity
	p.updatedAt = time.Now()
}

func (p *Part) Update(name, description string, unitPrice float64) error {
	if name == "" {
		return ErrPartNameRequired
	}
	if unitPrice <= 0 {
		return ErrPartInvalidPrice
	}
	p.name = name
	p.description = description
	p.unitPrice = unitPrice
	p.updatedAt = time.Now()
	return nil
}
