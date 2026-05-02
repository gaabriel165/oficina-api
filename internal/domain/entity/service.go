package entity

import (
	"time"

	"github.com/google/uuid"
)

type Service struct {
	id               string
	name             string
	description      string
	laborPrice       float64
	estimatedMinutes int
	createdAt        time.Time
	updatedAt        time.Time
}

func NewService(name, description string, laborPrice float64, estimatedMinutes int) (*Service, error) {
	if name == "" {
		return nil, ErrServiceNameRequired
	}
	if laborPrice <= 0 {
		return nil, ErrServiceInvalidLaborPrice
	}
	if estimatedMinutes <= 0 {
		return nil, ErrServiceInvalidEstimatedMin
	}

	now := time.Now()
	return &Service{
		id:               uuid.NewString(),
		name:             name,
		description:      description,
		laborPrice:       laborPrice,
		estimatedMinutes: estimatedMinutes,
		createdAt:        now,
		updatedAt:        now,
	}, nil
}

func RestoreService(id, name, description string, laborPrice float64, estimatedMinutes int, createdAt, updatedAt time.Time) *Service {
	return &Service{
		id:               id,
		name:             name,
		description:      description,
		laborPrice:       laborPrice,
		estimatedMinutes: estimatedMinutes,
		createdAt:        createdAt,
		updatedAt:        updatedAt,
	}
}

func (s *Service) ID() string               { return s.id }
func (s *Service) Name() string             { return s.name }
func (s *Service) Description() string      { return s.description }
func (s *Service) LaborPrice() float64      { return s.laborPrice }
func (s *Service) EstimatedMinutes() int    { return s.estimatedMinutes }
func (s *Service) CreatedAt() time.Time     { return s.createdAt }
func (s *Service) UpdatedAt() time.Time     { return s.updatedAt }

func (s *Service) Update(name, description string, laborPrice float64, estimatedMinutes int) error {
	if name == "" {
		return ErrServiceNameRequired
	}
	if laborPrice <= 0 {
		return ErrServiceInvalidLaborPrice
	}
	if estimatedMinutes <= 0 {
		return ErrServiceInvalidEstimatedMin
	}
	s.name = name
	s.description = description
	s.laborPrice = laborPrice
	s.estimatedMinutes = estimatedMinutes
	s.updatedAt = time.Now()
	return nil
}
