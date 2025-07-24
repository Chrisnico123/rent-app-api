package product

import (
	"context"
	"encoding/json"
	"fmt"
	"rent-application/domain"
	"rent-application/shared/helper"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductQuery interface {
	// Category operations
	CreateCategory(ctx context.Context, tx pgx.Tx, req domain.Category) error
	UpdateCategory(ctx context.Context, tx pgx.Tx, req domain.Category) error
	SoftDeleteCategory(ctx context.Context, tx pgx.Tx, id string) error

	GetAllCategories(ctx context.Context, db *pgxpool.Pool, filter domain.FilterSearchPagination) ([]domain.Category, error)
	CountAllCategories(ctx context.Context, db *pgxpool.Pool, filter domain.FilterSearchPagination) (int, error)
	GetCategoryByID(ctx context.Context, db *pgxpool.Pool, id string) (domain.Category, error)

	// Product operations
	CreateProduct(ctx context.Context, tx pgx.Tx, req domain.Product) error
	UpdateProduct(ctx context.Context, tx pgx.Tx, req domain.Product) error
	SoftDeleteProduct(ctx context.Context, tx pgx.Tx, id string) error
	UpdateAvailableProduct(ctx context.Context, tx pgx.Tx, id string, available bool) error

	GetProductByID(ctx context.Context, db *pgxpool.Pool, id string) (domain.ProductResponse, error)
	GetAllProducts(ctx context.Context, db *pgxpool.Pool, filter domain.FilterProduct) ([]domain.ProductResponseList, error)
	CountAllProducts(ctx context.Context, db *pgxpool.Pool, filter domain.FilterProduct) (int, error)
}

type ProductQueryImpl struct {
}

func NewProductQuery() ProductQuery {
	return &ProductQueryImpl{}
}

func (q *ProductQueryImpl) UpdateAvailableProduct(ctx context.Context, tx pgx.Tx, id string, available bool) error {
	updateQuery := `UPDATE product SET available = $1, updated_at = NOW() WHERE id = $2 AND delete_at IS NULL`
	result, err := tx.Exec(ctx, updateQuery, available, id)
	if err != nil {
		return err
	}

	// Check if any row was actually updated
	if result.RowsAffected() == 0 {
		return fmt.Errorf("product with id %s not found", id)
	}

	return nil
}

func (q *ProductQueryImpl) CountAllProducts(ctx context.Context, db *pgxpool.Pool, filter domain.FilterProduct) (int, error) {
	filterQuery, _, _ := filter.QueryBuildFilterProduct()
	query := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM product p
		WHERE delete_at IS NULL %s`, filterQuery)
	var count int
	err := db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (q *ProductQueryImpl) GetAllProducts(ctx context.Context, db *pgxpool.Pool, filter domain.FilterProduct) ([]domain.ProductResponseList, error) {
	filterQuery, sortQuery, paginationQuery := filter.QueryBuildFilterProduct()
	query := fmt.Sprintf(`
		SELECT 
			p.id, p.name, p.img, p.price, p.available, category.name AS category_name, p.created_at, p.updated_at
		FROM 
			product p
		LEFT JOIN category ON p.category_id = category.id
		WHERE p.delete_at IS NULL %s %s
		%s`, filterQuery, sortQuery, paginationQuery)

	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []domain.ProductResponseList
	for rows.Next() {
		var product domain.ProductResponseList
		var imgBytes []byte
		var createdAt, updatedAt interface{}
		if err := rows.Scan(&product.Id, &product.Name, &imgBytes, &product.Price, &product.Available, &product.Category, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		if len(imgBytes) > 0 {
			_ = json.Unmarshal(imgBytes, &product.Img)
		}
		// Convert time to Asia/Jakarta using helper
		if t, ok := createdAt.(time.Time); ok {
			product.CreatedAt = helper.ConvertToJakarta(t)
		}
		if t, ok := updatedAt.(time.Time); ok {
			product.UpdatedAt = helper.ConvertToJakarta(t)
		}
		products = append(products, product)
	}
	if len(products) == 0 {
		return nil, pgx.ErrNoRows
	}
	return products, nil
}

func (q *ProductQueryImpl) CreateProduct(ctx context.Context, tx pgx.Tx, req domain.Product) error {
	descJSON, err := json.Marshal(req.Description)
	if err != nil {
		return err
	}
	imgJSON, err := json.Marshal(req.Img)
	if err != nil {
		return err
	}
	// Convert []byte to string for JSONB input
	query := `INSERT INTO product (id, seller_id, name, category_id, price, description, img, available) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = tx.Exec(ctx, query, req.ID, req.SellerID, req.Name, req.CategoryID, req.Price, string(descJSON), string(imgJSON), req.Available)
	return err
}

func (q *ProductQueryImpl) UpdateProduct(ctx context.Context, tx pgx.Tx, req domain.Product) error {
	descJSON, err := json.Marshal(req.Description)
	if err != nil {
		return err
	}
	imgJSON, err := json.Marshal(req.Img)
	if err != nil {
		return err
	}

	query := `UPDATE product SET name = $1, category_id = $2, price = $3, description = $4, img = $5, available = $6 ,updated_at = NOW() WHERE id = $7`
	_, err = tx.Exec(ctx, query, req.Name, req.CategoryID, req.Price, string(descJSON), string(imgJSON), req.Available, req.ID)
	return err
}

func (q *ProductQueryImpl) SoftDeleteProduct(ctx context.Context, tx pgx.Tx, id string) error {
	query := `UPDATE product SET delete_at = NOW() WHERE id = $1`
	_, err := tx.Exec(ctx, query, id)
	return err
}

func (q *ProductQueryImpl) GetProductByID(ctx context.Context, db *pgxpool.Pool, id string) (domain.ProductResponse, error) {
	query := `
		SELECT 
			product.id, seller_id, product.name, category_id, price, product.available, description, product.img, product.created_at, product.updated_at, category.name AS category_name 
		FROM 
			product
		LEFT JOIN category ON product.category_id = category.id
		WHERE product.id = $1 AND product.delete_at IS NULL`

	var product domain.ProductResponse
	var descBytes, imgBytes []byte
	var createdAt, updatedAt interface{}
	err := db.QueryRow(ctx, query, id).Scan(&product.ID, &product.SellerID, &product.Name, &product.CategoryID, &product.Price, &product.Available, &descBytes, &imgBytes, &createdAt, &updatedAt, &product.CategoryName)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.ProductResponse{}, fmt.Errorf("product with id %s not found", id)
		}
		return domain.ProductResponse{}, err
	}
	// Unmarshal JSONB fields
	if len(descBytes) > 0 {
		_ = json.Unmarshal(descBytes, &product.Description)
	}
	if len(imgBytes) > 0 {
		_ = json.Unmarshal(imgBytes, &product.Img)
	}
	// Convert time to Asia/Jakarta using helper
	if t, ok := createdAt.(time.Time); ok {
		product.CreatedAt = helper.ConvertToJakarta(t)
	}
	if t, ok := updatedAt.(time.Time); ok {
		product.UpdatedAt = helper.ConvertToJakarta(t)
	}
	return product, nil
}

func (q *ProductQueryImpl) GetCategoryByID(ctx context.Context, db *pgxpool.Pool, id string) (domain.Category, error) {
	query := `SELECT id, name, img FROM category WHERE id = $1 AND delete_at IS NULL`
	var category domain.Category
	err := db.QueryRow(ctx, query, id).Scan(&category.ID, &category.Name, &category.Img)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.Category{}, fmt.Errorf("category with id %s not found", id)
		}
		return domain.Category{}, err
	}
	return category, nil
}

func (q *ProductQueryImpl) CreateCategory(ctx context.Context, tx pgx.Tx, req domain.Category) error {
	query := `INSERT INTO category (id, name, img) VALUES ($1, $2, $3)`
	_, err := tx.Exec(ctx, query, req.ID, req.Name, req.Img)
	return err
}

func (q *ProductQueryImpl) UpdateCategory(ctx context.Context, tx pgx.Tx, req domain.Category) error {
	query := `UPDATE category SET name = $1, img = $2 WHERE id = $3`
	_, err := tx.Exec(ctx, query, req.Name, req.Img, req.ID)
	return err
}

func (q *ProductQueryImpl) SoftDeleteCategory(ctx context.Context, tx pgx.Tx, id string) error {
	query := `UPDATE category SET delete_at = NOW() WHERE id = $1`
	_, err := tx.Exec(ctx, query, id)
	return err
}

func (q *ProductQueryImpl) GetAllCategories(ctx context.Context, db *pgxpool.Pool, filter domain.FilterSearchPagination) ([]domain.Category, error) {
	filterQuery, paginationQuery := filter.QueryBuildFilterSearchPagination()

	var categories []domain.Category
	query := fmt.Sprintf(`
		SELECT 
			c.id, c.name, c.img 
		FROM 
			category c
		WHERE 
			delete_at IS NULL %s
		%s
	`, filterQuery, paginationQuery)

	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category domain.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Img); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if len(categories) == 0 {
		return nil, pgx.ErrNoRows
	}

	return categories, nil
}

func (q *ProductQueryImpl) CountAllCategories(ctx context.Context, db *pgxpool.Pool, filter domain.FilterSearchPagination) (int, error) {
	filterQuery, _ := filter.QueryBuildFilterSearchPagination()

	query := fmt.Sprintf(`
		SELECT 
			COUNT(*) 
		FROM 
			category c
		WHERE 
			delete_at IS NULL %s`, filterQuery)

	var count int
	err := db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
