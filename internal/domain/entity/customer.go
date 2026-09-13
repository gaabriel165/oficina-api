package entity

import (
	"time"

	"github.com/gabrielcamargo/oficina-api/internal/domain/valueobject"
	"github.com/google/uuid"
)

type DocumentType string

const (
	DocumentTypeCPF  DocumentType = "CPF"
	DocumentTypeCNPJ DocumentType = "CNPJ"
)

type Customer struct {
	id           string
	name         string
	documentType DocumentType
	document     string
	phone        string
	email        string
	status       valueobject.CustomerStatus
	createdAt    time.Time
	updatedAt    time.Time
}

func NewCustomer(name, document, phone, email string) (*Customer, error) {
	if name == "" {
		return nil, ErrCustomerNameRequired
	}
	if phone == "" {
		return nil, ErrCustomerPhoneRequired
	}
	if email == "" {
		return nil, ErrCustomerEmailRequired
	}

	docType, docValue, err := parseDocument(document)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &Customer{
		id:           uuid.NewString(),
		name:         name,
		documentType: docType,
		document:     docValue,
		phone:        phone,
		email:        email,
		status:       valueobject.CustomerStatusActive,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

func RestoreCustomer(id, name, documentType, document, phone, email, status string, createdAt, updatedAt time.Time) *Customer {
	return &Customer{
		id:           id,
		name:         name,
		documentType: DocumentType(documentType),
		document:     document,
		phone:        phone,
		email:        email,
		status:       valueobject.CustomerStatus(status),
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

func (c *Customer) ID() string                         { return c.id }
func (c *Customer) Name() string                       { return c.name }
func (c *Customer) Document() string                   { return c.document }
func (c *Customer) DocumentType() DocumentType         { return c.documentType }
func (c *Customer) Phone() string                      { return c.phone }
func (c *Customer) Email() string                      { return c.email }
func (c *Customer) Status() valueobject.CustomerStatus { return c.status }
func (c *Customer) IsActive() bool                     { return c.status.IsActive() }
func (c *Customer) CreatedAt() time.Time               { return c.createdAt }
func (c *Customer) UpdatedAt() time.Time               { return c.updatedAt }

func (c *Customer) Update(name, phone, email string) error {
	if name == "" {
		return ErrCustomerNameRequired
	}
	if phone == "" {
		return ErrCustomerPhoneRequired
	}
	if email == "" {
		return ErrCustomerEmailRequired
	}
	c.name = name
	c.phone = phone
	c.email = email
	c.updatedAt = time.Now()
	return nil
}

func parseDocument(document string) (DocumentType, string, error) {
	cpf, err := valueobject.NewCPF(document)
	if err == nil {
		return DocumentTypeCPF, cpf.Value(), nil
	}

	cnpj, err := valueobject.NewCNPJ(document)
	if err == nil {
		return DocumentTypeCNPJ, cnpj.Value(), nil
	}

	return "", "", ErrInvalidDocument
}

func (c *Customer) Activate() {
	c.status = valueobject.CustomerStatusActive
	c.updatedAt = time.Now()
}

func (c *Customer) Deactivate() {
	c.status = valueobject.CustomerStatusInactive
	c.updatedAt = time.Now()
}
