package product

import "fmt"

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Img  string `json:"img"`
}

func (c Category) ValidateCategory() error {
	if c.Name == "" {
		return fmt.Errorf("category name cannot be empty")
	}
	if c.Img == "" {
		return fmt.Errorf("category image cannot be empty")
	}
	return nil
}

type Product struct {
	ID          string   `json:"id"`
	SellerID    string   `json:"seller_id"`
	Name        string   `json:"name"`
	Description []string `json:"description"`
	Available   bool     `json:"available"`
	Price       float64  `json:"price"`
	CategoryID  string   `json:"category_id"`
	Img         []string `json:"img"`
}

func (p Product) ValidateProduct() error {
	if p.Name == "" {
		return fmt.Errorf("product name cannot be empty")
	}
	if len(p.Description) == 0 {
		return fmt.Errorf("product description cannot be empty")
	}
	if p.Price <= 0 {
		return fmt.Errorf("product price must be greater than zero")
	}
	if p.CategoryID == "" {
		return fmt.Errorf("product category ID cannot be empty")
	}
	if len(p.Img) == 0 {
		return fmt.Errorf("product images cannot be empty")
	}
	return nil
}

type BannerRequest struct {
	Name string `json:"name"`
	Img  string `json:"img"`
}

func (c BannerRequest) ValidateBanner() error {
	if c.Name == "" {
		return fmt.Errorf("banner name cannot be empty")
	}
	if c.Img == "" {
		return fmt.Errorf("banner image cannot be empty")
	}
	return nil
}
