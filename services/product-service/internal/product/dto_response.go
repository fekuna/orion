package product

import (
	"time"

	"github.com/google/uuid"
)

// ProductResponse is the outbound representation of a single product.
// Pagination metadata is handled by response.Meta in the standard envelope.
type ProductResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CategoryID  uuid.UUID `json:"category_id"`
	Available   bool      `json:"available"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
