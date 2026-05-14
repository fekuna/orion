package product

import (
	"context"
	"fmt"

	"github.com/fekuna/orion-v2/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type useCaseImpl struct {
	repo Repository
}

// NewUseCase creates a new product UseCase backed by the given Repository.
// No *zap.Logger is accepted — logging uses the request-scoped logger from ctx
// (injected by the server's loggerMiddleware) so every log entry automatically
// carries request_id, method, and URI without any extra work.
func NewUseCase(repo Repository) UseCase {
	return &useCaseImpl{repo: repo}
}

// GetProducts — read operation, no usecase-level logging.
// The HTTP access log already records every request; logging each read would be noise.
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

// GetProductByID — read operation, no usecase-level logging.
// ErrNotFound is a normal business outcome, not an anomaly worth logging.
func (uc *useCaseImpl) GetProductByID(ctx context.Context, id uuid.UUID) (*Product, error) {
	p, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get product by id: %w", err)
	}
	return p, nil
}

// CreateProduct — logs a business event at Info on success.
// Errors bubble up unwrapped; the handler layer logs them with HTTP context.
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

	logger.FromContext(ctx).Info("product created",
		zap.String("product_id", created.ID.String()),
		zap.String("name", created.Name),
		zap.Float64("price", created.Price),
	)

	return created, nil
}

// UpdateProduct — logs the mutation at Info and warns if stock hits zero.
// Stock depletion is a business anomaly worth alerting on.
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

	log := logger.FromContext(ctx)
	log.Info("product updated", zap.String("product_id", updated.ID.String()))

	// Business anomaly: alert when stock reaches zero so ops can restock.
	if updated.Stock == 0 {
		log.Warn("product stock depleted",
			zap.String("product_id", updated.ID.String()),
			zap.String("name", updated.Name),
		)
	}

	return updated, nil
}

// DeleteProduct — logs the deletion at Info on success.
// Deletes are high-value audit events — always worth recording.
func (uc *useCaseImpl) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete product: %w", err)
	}

	logger.FromContext(ctx).Info("product deleted",
		zap.String("product_id", id.String()),
	)

	return nil
}
