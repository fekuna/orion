package product

import (
	"context"

	"github.com/google/uuid"
)

// UseCase defines the application business logic for the product domain.
type UseCase interface {
	GetProducts(ctx context.Context, filter Filter) ([]*Product, int, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*Product, error)
	CreateProduct(ctx context.Context, req CreateProductRequest) (*Product, error)
	UpdateProduct(ctx context.Context, id uuid.UUID, req UpdateProductRequest) (*Product, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
}
