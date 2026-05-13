package product

// toResponse converts a domain Product entity to its HTTP response representation.
func toResponse(p *Product) *ProductResponse {
	return &ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		CategoryID:  p.CategoryID,
		Available:   p.IsAvailable(),
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

// toListData converts a slice of domain Products to a slice of response DTOs.
func toListData(products []*Product) []*ProductResponse {
	data := make([]*ProductResponse, len(products))
	for i, p := range products {
		data[i] = toResponse(p)
	}
	return data
}
