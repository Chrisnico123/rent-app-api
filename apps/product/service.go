package product

import (
	"context"
	"errors"
	"rent-application/domain"
	"rent-application/shared/helper"
	"rent-application/shared/web"

	"github.com/jackc/pgx/v5"
)

type ProductService interface {
	// Category operations
	CreateCategory(ctx context.Context, category Category) error
	UpdateCategory(ctx context.Context, category Category) error
	SoftDeleteCategory(ctx context.Context, id string) error

	GetAllCategories(ctx context.Context, filter web.FilterSearchPagination) ([]web.CategoryResponse, int, error)
	GetCategoryByID(ctx context.Context, id string) (web.CategoryResponse, error)

	// Product operations
	CreateProduct(ctx context.Context, product Product) error
	UpdateProduct(ctx context.Context, product Product) error
	SoftDeleteProduct(ctx context.Context, id string) error

	GetProductByID(ctx context.Context, id string) (web.ProductResponse, error)
	GetAllProducts(ctx context.Context, filter web.FilterProduct) ([]web.ProductResponseList, int, error)
	GetProductBookById(ctx context.Context, id, user_id string) (web.GetDataProductBook, error)

	// Banner operations
	CreateBanner(ctx context.Context, banner BannerRequest) error
	UpdateBanner(ctx context.Context, id string, banner BannerRequest) error
	DeleteBanner(ctx context.Context, id string) error

	GetBannerById(ctx context.Context, id string) (web.Banner, error)
	GetListBanner(ctx context.Context, filter web.FilterBannerPagination) ([]web.Banner, int, error)
}

type productService struct {
	productRepository ProductRepository
}

func NewProductService(productRepo ProductRepository) ProductService {
	return &productService{
		productRepository: productRepo,
	}
}

// GetProductBookById implements ProductService.
func (s *productService) GetProductBookById(ctx context.Context, id string, user_id string) (web.GetDataProductBook, error) {
	data, _ := s.productRepository.GetProductByID(ctx, id)
	if data.ID == "" {
		return web.GetDataProductBook{}, web.ErrNotFound("data not found")
	}

	product, err := s.productRepository.GetProductBookById(ctx, domain.GetProductBook{
		ProductId: id,
		UserId:    user_id,
	})
	if err != nil {
		return web.GetDataProductBook{}, err
	}

	return web.GetDataProductBook(product), nil
}

func (s *productService) CreateBanner(ctx context.Context, banner BannerRequest) error {
	if err := banner.ValidateBanner(); err != nil {
		return err
	}

	domainBanner := domain.Banner{
		Id:   helper.GenerateId(),
		Name: banner.Name,
		Img:  banner.Img,
	}

	return s.productRepository.CreateBanner(ctx, domainBanner)
}

func (s *productService) UpdateBanner(ctx context.Context, id string, banner BannerRequest) error {
	if err := banner.ValidateBanner(); err != nil {
		return err
	}

	data, _ := s.productRepository.GetBannerById(ctx, id)
	if data.Id == "" {
		return web.ErrNotFound("banner not found")
	}

	domainBanner := domain.Banner{
		Id:   id,
		Name: banner.Name,
		Img:  banner.Img,
	}

	return s.productRepository.UpdateBanner(ctx, domainBanner)
}

func (s *productService) DeleteBanner(ctx context.Context, id string) error {
	if id == "" {
		return web.ErrBadRequest("banner ID cannot be empty")
	}

	data, _ := s.productRepository.GetBannerById(ctx, id)
	if data.Id == "" {
		return web.ErrNotFound("banner not found")
	}

	return s.productRepository.DeleteBanner(ctx, id)
}

func (s *productService) GetBannerById(ctx context.Context, id string) (web.Banner, error) {
	data, err := s.productRepository.GetBannerById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return web.Banner{}, web.ErrNotFound("no active banner found")
		}
		return web.Banner{}, err
	}

	return web.Banner{
		Id:   data.Id,
		Name: data.Name,
		Img:  data.Img,
	}, nil
}

func (s *productService) GetListBanner(ctx context.Context, filter web.FilterBannerPagination) ([]web.Banner, int, error) {
	banners, total, err := s.productRepository.GetListBanner(ctx, domain.ToDomainFilterBannerPagination(filter))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []web.Banner{}, 0, nil
		}
		return nil, 0, err
	}

	var bannerResponses []web.Banner
	for _, banner := range banners {
		bannerResponses = append(bannerResponses, web.Banner{
			Id:   banner.Id,
			Name: banner.Name,
			Img:  banner.Img,
		})
	}

	return bannerResponses, total, nil
}

func (s *productService) GetAllProducts(ctx context.Context, filter web.FilterProduct) ([]web.ProductResponseList, int, error) {
	domainFilter := domain.ToDomainFilterProduct(filter)

	products, total, err := s.productRepository.GetAllProducts(ctx, domainFilter)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []web.ProductResponseList{}, 0, nil
		}
		return nil, 0, err
	}

	var productResponses []web.ProductResponseList
	for _, prod := range products {
		productResponses = append(productResponses, web.ProductResponseList{
			ID:        prod.Id,
			Name:      prod.Name,
			Img:       prod.Img[0],
			Category:  prod.Category,
			Price:     prod.Price,
			Available: prod.Available,
			CreatedAt: prod.CreatedAt,
			UpdatedAt: prod.UpdatedAt,
		})
	}
	return productResponses, total, nil
}

func (s *productService) CreateProduct(ctx context.Context, product Product) error {
	if err := product.ValidateProduct(); err != nil {
		return err
	}

	domainProduct := domain.Product{
		ID:          helper.GenerateId(),
		SellerID:    product.SellerID,
		Available:   product.Available,
		Name:        product.Name,
		CategoryID:  product.CategoryID,
		Price:       product.Price,
		Description: product.Description,
		Img:         product.Img,
	}

	return s.productRepository.CreateProduct(ctx, domainProduct)
}

func (s *productService) UpdateProduct(ctx context.Context, product Product) error {
	if err := product.ValidateProduct(); err != nil {
		return err
	}

	data, err := s.productRepository.GetProductByID(ctx, product.ID)
	if err != nil {
		return err
	}
	if data.ID == "" {
		return web.ErrNotFound("product not found")
	}

	domainProduct := domain.Product{
		ID:          product.ID,
		SellerID:    data.SellerID, // Keep the existing seller ID
		Name:        product.Name,
		Available:   product.Available,
		CategoryID:  product.CategoryID,
		Price:       product.Price,
		Description: product.Description,
		Img:         product.Img,
	}

	return s.productRepository.UpdateProduct(ctx, domainProduct)
}

func (s *productService) SoftDeleteProduct(ctx context.Context, id string) error {
	if id == "" {
		return web.ErrBadRequest("product ID cannot be empty")
	}

	data, err := s.productRepository.GetProductByID(ctx, id)
	if err != nil {
		return err
	}
	if data.ID == "" {
		return web.ErrNotFound("product not found")
	}

	return s.productRepository.SoftDeleteProduct(ctx, id)
}

func (s *productService) GetProductByID(ctx context.Context, id string) (web.ProductResponse, error) {
	if id == "" {
		return web.ProductResponse{}, web.ErrBadRequest("product ID cannot be empty")
	}

	data, err := s.productRepository.GetProductByID(ctx, id)
	if err != nil {
		return web.ProductResponse{}, err
	}
	if data.ID == "" {
		return web.ProductResponse{}, web.ErrNotFound("product not found")
	}

	return web.ProductResponse{
		ID:           data.ID,
		Name:         data.Name,
		Description:  data.Description,
		Available:    data.Available,
		Img:          data.Img,
		CategoryID:   data.CategoryID,
		CategoryName: data.CategoryName,
		Price:        data.Price,
		SellerID:     data.SellerID,
		CreatedAt:    data.CreatedAt,
		UpdatedAt:    data.UpdatedAt,
	}, nil
}

func (s *productService) GetCategoryByID(ctx context.Context, id string) (web.CategoryResponse, error) {
	if id == "" {
		return web.CategoryResponse{}, web.ErrBadRequest("category ID cannot be empty")
	}

	data, err := s.productRepository.GetCategoryByID(ctx, id)
	if err != nil {
		return web.CategoryResponse{}, err
	}
	if data.ID == "" {
		return web.CategoryResponse{}, web.ErrNotFound("category not found")
	}

	return web.CategoryResponse{
		ID:   data.ID,
		Name: data.Name,
		Img:  data.Img,
	}, nil
}

func (s *productService) CreateCategory(ctx context.Context, category Category) error {
	if err := category.ValidateCategory(); err != nil {
		return err
	}

	domainCategory := domain.Category{
		ID:   helper.GenerateId(),
		Name: category.Name,
		Img:  category.Img,
	}

	return s.productRepository.CreateCategory(ctx, domainCategory)
}

func (s *productService) UpdateCategory(ctx context.Context, category Category) error {
	if err := category.ValidateCategory(); err != nil {
		return err
	}

	data, err := s.productRepository.GetCategoryByID(ctx, category.ID)
	if err != nil {
		return err
	}
	if data.ID == "" {
		return web.ErrNotFound("category not found")
	}

	domainCategory := domain.Category{
		ID:   category.ID,
		Name: category.Name,
		Img:  category.Img,
	}

	return s.productRepository.UpdateCategory(ctx, domainCategory)
}

func (s *productService) SoftDeleteCategory(ctx context.Context, id string) error {
	if id == "" {
		return web.ErrBadRequest("category ID cannot be empty")
	}

	data, err := s.productRepository.GetCategoryByID(ctx, id)
	if err != nil {
		return err
	}
	if data.ID == "" {
		return web.ErrNotFound("category not found")
	}

	return s.productRepository.SoftDeleteCategory(ctx, id)
}

func (s *productService) GetAllCategories(ctx context.Context, filter web.FilterSearchPagination) ([]web.CategoryResponse, int, error) {
	domainFilter := domain.ToDomainFilterSearchPagination(filter)

	categories, total, err := s.productRepository.GetAllCategories(ctx, domainFilter)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []web.CategoryResponse{}, 0, nil
		}
		return nil, 0, err
	}

	var categoryResponses []web.CategoryResponse
	for _, cat := range categories {
		categoryResponses = append(categoryResponses, web.CategoryResponse{
			ID:   cat.ID,
			Name: cat.Name,
			Img:  cat.Img,
		})
	}
	return categoryResponses, total, nil
}
