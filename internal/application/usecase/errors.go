package usecase

import "errors"

var (
	ErrDocumentAlreadyRegistered = errors.New("document already registered")
	ErrPlateAlreadyRegistered    = errors.New("plate already registered")
	ErrEmailAlreadyRegistered    = errors.New("email already registered")
	ErrCNPJNotActive             = errors.New("CNPJ is not active")
	ErrInvalidCredentials        = errors.New("invalid credentials")
	ErrForbidden                 = errors.New("access to this resource is not allowed")
)
