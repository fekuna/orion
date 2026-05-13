package product

import "errors"

// Sentinel errors for the product domain.
// Handlers map these to appropriate HTTP status codes.
var (
	ErrNotFound      = errors.New("product not found")
	ErrAlreadyExists = errors.New("product already exists")
	ErrInvalidInput  = errors.New("invalid product input")
)
