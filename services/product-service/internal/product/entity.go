package product

import (
	"time"

	"github.com/google/uuid"
)

// Product is the core domain entity.
type Product struct {
	ID          uuid.UUID
	Name        string
	Description string
	Price       float64
	Stock       int
	CategoryID  uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// IsAvailable returns true when there is stock remaining.
func (p *Product) IsAvailable() bool {
	return p.Stock > 0
}
