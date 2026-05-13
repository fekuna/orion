package product

import "github.com/google/uuid"

// CreateProductRequest is the payload for creating a new product.
type CreateProductRequest struct {
	Name        string    `json:"name"        validate:"required,min=1,max=255"`
	Description string    `json:"description" validate:"max=2000"`
	Price       float64   `json:"price"       validate:"required,gt=0"`
	Stock       int       `json:"stock"       validate:"gte=0"`
	CategoryID  uuid.UUID `json:"category_id" validate:"required"`
}

// UpdateProductRequest uses pointers so only provided fields are updated (partial update).
type UpdateProductRequest struct {
	Name        *string  `json:"name"        validate:"omitempty,min=1,max=255"`
	Description *string  `json:"description" validate:"omitempty,max=2000"`
	Price       *float64 `json:"price"       validate:"omitempty,gt=0"`
	Stock       *int     `json:"stock"       validate:"omitempty,gte=0"`
}
