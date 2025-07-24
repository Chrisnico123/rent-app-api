package product

import (
	"context"
	"rent-application/domain"
	"rent-application/internal/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository interface {
	// Category operations
	CreateCategory(ctx context.Context, category domain.Category) error
	UpdateCategory(ctx context.Context, category domain.Category) error
	SoftDeleteCategory(ctx context.Context, id string) error

	GetCategoryByID(ctx context.Context, id string) (domain.Category, error)
	GetAllCategories(ctx context.Context, filter domain.FilterSearchPagination) ([]domain.Category, int, error)

	// Product operations
	CreateProduct(ctx context.Context, product domain.Product) error
	UpdateProduct(ctx context.Context, product domain.Product) error
	SoftDeleteProduct(ctx context.Context, id string) error

	GetProductByID(ctx context.Context, id string) (domain.ProductResponse, error)
	GetAllProducts(ctx context.Context, filter domain.FilterProduct) ([]domain.ProductResponseList, int, error)
}

type productRepository struct {
	db           database.Store
	productQuery ProductQuery
}

func NewProductRepository(db database.Store, productQuery ProductQuery) ProductRepository {
	return &productRepository{
		db:           db,
		productQuery: productQuery,
	}
}

func (r *productRepository) GetAllProducts(ctx context.Context, filter domain.FilterProduct) ([]domain.ProductResponseList, int, error) {
	var products []domain.ProductResponseList
	var total int

	err := r.db.WithoutTransaction(ctx, func(db *pgxpool.Pool) error {
		var err error
		if products, err = r.productQuery.GetAllProducts(ctx, db, filter); err != nil {
			return err
		}

		if total, err = r.productQuery.CountAllProducts(ctx, db, filter); err != nil {
			return err
		}

		return nil
	})

	return products, total, err
}

func (r *productRepository) CreateProduct(ctx context.Context, product domain.Product) error {
	return r.db.WithTransaction(ctx, func(tx pgx.Tx) error {
		return r.productQuery.CreateProduct(ctx, tx, product)
	})
}

func (r *productRepository) UpdateProduct(ctx context.Context, product domain.Product) error {
	return r.db.WithTransaction(ctx, func(tx pgx.Tx) error {
		return r.productQuery.UpdateProduct(ctx, tx, product)
	})
}

func (r *productRepository) SoftDeleteProduct(ctx context.Context, id string) error {
	return r.db.WithTransaction(ctx, func(tx pgx.Tx) error {
		return r.productQuery.SoftDeleteProduct(ctx, tx, id)
	})
}

func (r *productRepository) GetProductByID(ctx context.Context, id string) (domain.ProductResponse, error) {
	var product domain.ProductResponse
	err := r.db.WithoutTransaction(ctx, func(db *pgxpool.Pool) error {
		var err error
		product, err = r.productQuery.GetProductByID(ctx, db, id)
		return err
	})
	if err != nil {
		return domain.ProductResponse{}, err
	}
	return product, nil
}

func (r *productRepository) GetCategoryByID(ctx context.Context, id string) (domain.Category, error) {
	var category domain.Category
	var err error

	err = r.db.WithoutTransaction(ctx, func(db *pgxpool.Pool) error {
		category, err = r.productQuery.GetCategoryByID(ctx, db, id)
		if err != nil {
			return err
		}
		return nil
	})

	return category, err
}

func (r *productRepository) CreateCategory(ctx context.Context, category domain.Category) error {
	return r.db.WithTransaction(ctx, func(tx pgx.Tx) error {
		return r.productQuery.CreateCategory(ctx, tx, category)
	})
}

func (r *productRepository) UpdateCategory(ctx context.Context, category domain.Category) error {
	return r.db.WithTransaction(ctx, func(tx pgx.Tx) error {
		return r.productQuery.UpdateCategory(ctx, tx, category)
	})
}

func (r *productRepository) SoftDeleteCategory(ctx context.Context, id string) error {
	return r.db.WithTransaction(ctx, func(tx pgx.Tx) error {
		return r.productQuery.SoftDeleteCategory(ctx, tx, id)
	})
}

func (r *productRepository) GetAllCategories(ctx context.Context, filter domain.FilterSearchPagination) ([]domain.Category, int, error) {
	var categories []domain.Category
	var total int

	err := r.db.WithoutTransaction(ctx, func(db *pgxpool.Pool) error {
		var err error
		if categories, err = r.productQuery.GetAllCategories(ctx, db, filter); err != nil {
			return err
		}

		if total, err = r.productQuery.CountAllCategories(ctx, db, filter); err != nil {
			return err
		}

		return nil
	})

	return categories, total, err
}
