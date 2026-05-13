package product

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type useCaseImpl struct {
	repo Repository
}

// NewUseCase creates a new product UseCase backed by the given Repository.
// Logging is intentionally omitted here — it belongs at the handler boundary
// where HTTP context (request ID, method, URI) is available.
func NewUseCase(repo Repository) UseCase {
	return &useCaseImpl{repo: repo}
}

func (uc *useCaseImpl) GetProducts(ctx context.Context, filter Filter) ([]*Product, int, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}

	products, total, err := uc.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("get products: %w", err)
	}
	return products, total, nil
}

func (uc *useCaseImpl) GetProductByID(ctx context.Context, id uuid.UUID) (*Product, error) {
	p, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get product by id: %w", err)
	}
	return p, nil
}

func (uc *useCaseImpl) CreateProduct(ctx context.Context, req CreateProductRequest) (*Product, error) {
	p := &Product{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		CategoryID:  req.CategoryID,
	}

	created, err := uc.repo.Create(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("create product: %w", err)
	}
	return created, nil
}

func (uc *useCaseImpl) UpdateProduct(ctx context.Context, id uuid.UUID, req UpdateProductRequest) (*Product, error) {
	existing, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("update product — find: %w", err)
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.Price != nil {
		existing.Price = *req.Price
	}
	if req.Stock != nil {
		existing.Stock = *req.Stock
	}

	updated, err := uc.repo.Update(ctx, existing)
	if err != nil {
		return nil, fmt.Errorf("update product: %w", err)
	}
	return updated, nil
}

func (uc *useCaseImpl) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete product: %w", err)
	}
	return nil
}
