package product

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// productColumns is the canonical SELECT / RETURNING column list for the
// products table. Defined once so adding/removing a column is a single edit.
const productColumns = `id, name, description, price, stock, category_id, created_at, updated_at`

type postgresRepo struct {
	db *pgxpool.Pool
}

// NewPostgresRepo creates a PostgreSQL-backed product Repository.
func NewPostgresRepo(db *pgxpool.Pool) Repository {
	return &postgresRepo{db: db}
}

// scanProduct reads a single product row from any pgx.Row-compatible source.
// Centralises all Scan calls so adding a column is a single edit here.
func scanProduct(row interface {
	Scan(dest ...any) error
}) (*Product, error) {
	p := &Product{}
	err := row.Scan(
		&p.ID, &p.Name, &p.Description, &p.Price,
		&p.Stock, &p.CategoryID, &p.CreatedAt, &p.UpdatedAt,
	)
	return p, err
}

func (r *postgresRepo) FindAll(ctx context.Context, filter Filter) ([]*Product, int, error) {
	var total int
	countQ := `SELECT COUNT(*) FROM products WHERE ($1 = '' OR name ILIKE '%' || $1 || '%')`
	if err := r.db.QueryRow(ctx, countQ, filter.Name).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("postgres: count products: %w", err)
	}

	offset := (filter.Page - 1) * filter.Limit
	q := fmt.Sprintf(`
		SELECT %s FROM products
		WHERE ($1 = '' OR name ILIKE '%%' || $1 || '%%')
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`, productColumns)

	rows, err := r.db.Query(ctx, q, filter.Name, filter.Limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("postgres: query products: %w", err)
	}
	defer rows.Close()

	var products []*Product
	for rows.Next() {
		p, err := scanProduct(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("postgres: scan product: %w", err)
		}
		products = append(products, p)
	}

	return products, total, nil
}

func (r *postgresRepo) FindByID(ctx context.Context, id uuid.UUID) (*Product, error) {
	q := fmt.Sprintf(`SELECT %s FROM products WHERE id = $1`, productColumns)

	p, err := scanProduct(r.db.QueryRow(ctx, q, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("postgres: find product by id: %w", err)
	}
	return p, nil
}

func (r *postgresRepo) Create(ctx context.Context, p *Product) (*Product, error) {
	q := fmt.Sprintf(`
		INSERT INTO products (id, name, description, price, stock, category_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING %s`, productColumns)

	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now

	created, err := scanProduct(r.db.QueryRow(ctx, q,
		p.ID, p.Name, p.Description, p.Price,
		p.Stock, p.CategoryID, p.CreatedAt, p.UpdatedAt,
	))
	if err != nil {
		return nil, fmt.Errorf("postgres: create product: %w", err)
	}
	return created, nil
}

func (r *postgresRepo) Update(ctx context.Context, p *Product) (*Product, error) {
	q := fmt.Sprintf(`
		UPDATE products
		SET name=$1, description=$2, price=$3, stock=$4, updated_at=$5
		WHERE id=$6
		RETURNING %s`, productColumns)

	p.UpdatedAt = time.Now().UTC()

	updated, err := scanProduct(r.db.QueryRow(ctx, q,
		p.Name, p.Description, p.Price, p.Stock, p.UpdatedAt, p.ID,
	))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("postgres: update product: %w", err)
	}
	return updated, nil
}

func (r *postgresRepo) Delete(ctx context.Context, id uuid.UUID) error {
	q := `DELETE FROM products WHERE id = $1`
	tag, err := r.db.Exec(ctx, q, id)
	if err != nil {
		return fmt.Errorf("postgres: delete product: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
