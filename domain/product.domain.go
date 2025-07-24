package domain

import (
	"fmt"
	"rent-application/shared/helper"
	"rent-application/shared/web"
	"strconv"
)

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

type ProductResponseList struct {
	Id        string
	Name      string
	Img       []string
	Price     float64
	Available bool
	Category  string
	CreatedAt string
	UpdatedAt string
}

type ProductResponse struct {
	ID           string   `json:"id"`
	SellerID     string   `json:"seller_id"`
	Name         string   `json:"name"`
	Description  []string `json:"description"`
	Available    bool     `json:"available"`
	Price        float64  `json:"price"`
	CategoryID   string   `json:"category_id"`
	CategoryName string   `json:"category_name"`
	Img          []string `json:"img"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
}

type FilterProduct struct {
	Search     string
	CategoryID string
	Sort       string
	Pagination Pagination
}

func ToDomainFilterProduct(q web.FilterProduct) FilterProduct {
	pageInt, _ := strconv.Atoi(q.Page)
	limitInt, _ := strconv.Atoi(q.Limit)

	var page, limit string

	if q.Page != "" {
		if q.Limit != "" {
			page = fmt.Sprintf("%d", (pageInt-1)*limitInt)
			limit = fmt.Sprintf("%d", limitInt)
		} else {
			page = fmt.Sprintf("%d", (pageInt-1)*5)
			limit = fmt.Sprintf("%d", 5)
		}
	} else {
		if q.Limit != "" {
			page = "1"
			limit = fmt.Sprintf("%d", limitInt)
		}
	}

	return FilterProduct{
		Search:     q.Search,
		CategoryID: q.CategoryID,
		Sort:       q.Sort,
		Pagination: Pagination{
			Page:  page,
			Limit: limit,
		},
	}
}

func (f *FilterProduct) QueryBuildFilterProduct() (filter, sort, pagination string) {
	filter = BuildProductQuerySearch(f.Search)

	if f.CategoryID != "" {
		filter += fmt.Sprintf(" AND category_id = '%s'", f.CategoryID)
	}

	if f.Sort != "" {
		switch f.Sort {
		case "1":
			sort += " ORDER BY p.price ASC"
		case "2":
			sort += " ORDER BY p.price DESC"
		}
	}

	if f.Pagination.Page != "" && f.Pagination.Limit != "" {
		pagination = fmt.Sprintf("OFFSET %s LIMIT %s", f.Pagination.Page, f.Pagination.Limit)
	}
	return
}

func BuildProductQuerySearch(search string) (filter string) {
	fieldsAndValues := map[string]string{
		"p.name": search,
	}

	filter = helper.BuildManyColumnFilter(fieldsAndValues, filter)

	return
}
