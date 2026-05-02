package repository

import "errors"

var (
	ErrCustomerNotFound     = errors.New("customer not found")
	ErrVehicleNotFound      = errors.New("vehicle not found")
	ErrPartNotFound         = errors.New("part not found")
	ErrServiceNotFound      = errors.New("service not found")
	ErrServiceOrderNotFound = errors.New("service order not found")
	ErrUserNotFound         = errors.New("user not found")
)
