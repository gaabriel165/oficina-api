package model

import (
	"time"

	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
)

type CustomerModel struct {
	ID           string    `gorm:"type:uuid;primaryKey;column:id"`
	Name         string    `gorm:"type:varchar(255);not null;column:name"`
	Document     string    `gorm:"type:varchar(14);not null;uniqueIndex;column:document"`
	DocumentType string    `gorm:"type:varchar(4);not null;column:document_type"`
	Phone        string    `gorm:"type:varchar(20);not null;column:phone"`
	Email        string    `gorm:"type:varchar(255);not null;column:email"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (CustomerModel) TableName() string { return "customers" }

func (m *CustomerModel) ToDomain() *entity.Customer {
	return entity.RestoreCustomer(m.ID, m.Name, m.DocumentType, m.Document, m.Phone, m.Email, m.CreatedAt, m.UpdatedAt)
}

func CustomerModelFromDomain(c *entity.Customer) *CustomerModel {
	return &CustomerModel{
		ID:           c.ID(),
		Name:         c.Name(),
		Document:     c.Document(),
		DocumentType: string(c.DocumentType()),
		Phone:        c.Phone(),
		Email:        c.Email(),
		CreatedAt:    c.CreatedAt(),
		UpdatedAt:    c.UpdatedAt(),
	}
}
