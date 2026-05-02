package entity

import "errors"

var (
	ErrCustomerNameRequired  = errors.New("customer name is required")
	ErrCustomerPhoneRequired = errors.New("customer phone is required")
	ErrCustomerEmailRequired = errors.New("customer email is required")
	ErrInvalidDocument       = errors.New("document must be a valid CPF or CNPJ")

	ErrVehicleBrandRequired    = errors.New("vehicle brand is required")
	ErrVehicleModelRequired    = errors.New("vehicle model is required")
	ErrVehicleInvalidYear      = errors.New("vehicle year is invalid")
	ErrVehicleCustomerRequired = errors.New("vehicle customer id is required")

	ErrPartNameRequired          = errors.New("part name is required")
	ErrPartInvalidPrice          = errors.New("part unit price must be greater than zero")
	ErrPartInvalidStock          = errors.New("part stock quantity cannot be negative")
	ErrPartInsufficientStock     = errors.New("insufficient stock quantity")

	ErrServiceNameRequired        = errors.New("service name is required")
	ErrServiceInvalidLaborPrice   = errors.New("service labor price must be greater than zero")
	ErrServiceInvalidEstimatedMin = errors.New("service estimated minutes must be greater than zero")

	ErrServiceOrderCustomerRequired      = errors.New("service order customer id is required")
	ErrServiceOrderVehicleRequired       = errors.New("service order vehicle id is required")
	ErrServiceOrderVehicleNotOwnedByCustomer = errors.New("vehicle does not belong to the specified customer")
	ErrServiceOrderInvalidTransition     = errors.New("service order status transition not allowed")
	ErrServiceOrderNoItems               = errors.New("service order must have at least one service or part")
	ErrServiceOrderItemInvalidQty        = errors.New("service order part quantity must be greater than zero")
)
