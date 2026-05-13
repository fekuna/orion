package product

import (
	"context"

	"github.com/google/uuid"
)

// Filter holds optional query parameters for listing products.
type Filter struct {
	Page  int
	Limit int
	Name  string
}

// Repository defines the persistence contract for the product domain.
// All implementations (e.g. postgresRepo) live in this package.
type Repository interface {
	FindAll(ctx context.Context, filter Filter) ([]*Product, int, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Product, error)
	Create(ctx context.Context, p *Product) (*Product, error)
	Update(ctx context.Context, p *Product) (*Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
