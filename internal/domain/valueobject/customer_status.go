package valueobject

import "errors"

var ErrInvalidCustomerStatus = errors.New("customer status must be active or inactive")

type CustomerStatus string

const (
	CustomerStatusActive   CustomerStatus = "active"
	CustomerStatusInactive CustomerStatus = "inactive"
)

func NewCustomerStatus(value string) (CustomerStatus, error) {
	status := CustomerStatus(value)
	if status != CustomerStatusActive && status != CustomerStatusInactive {
		return "", ErrInvalidCustomerStatus
	}
	return status, nil
}

func (s CustomerStatus) IsActive() bool {
	return s == CustomerStatusActive
}

func (s CustomerStatus) String() string {
	return string(s)
}
