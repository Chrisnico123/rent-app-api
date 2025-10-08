package domain

import (
	"fmt"
	"rent-application/shared/helper"
	"rent-application/shared/web"
	"strconv"
	"time"
)

type OrderUser struct {
	ID        string
	UserID    string
	PaymentId string
	OrderID   string
	Type      string
	StartDate string
	EndDate   string
	Desc      string
	ProductID *string
	Price     float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// PaymentHistory domain struct
type PaymentHistory struct {
	ID        string
	UserID    string
	PaymentId string
	OrderID   string
	VAID      string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiredAt *time.Time
}

type PaymentResponse struct {
	ID          string
	Username    string
	Email       string
	UserID      string
	Method      string
	PaymentId   string
	OrderID     string
	ProductId   string
	StartDate   string
	EndDate     string
	BookingDays uint16
	VAID        string
	Status      string
	CreatedAt   string
	UpdatedAt   string
	ExpiredAt   string
}

type PaymentHistoryResponse struct {
	OrderId      string
	CustomerName string
	Email        string
	StartDate    string
	EndDate      string
	ProductName  string
	Price        string
	Img          []string
	Status       string
	Method       string
	CreatedDate  string
	ExpiredDate  string
}
type FilterPaymentHistory struct {
	Search     string
	UserId     string
	Status     string
	Level      string
	Pagination Pagination
}

func ToDomainFilterPaymentHistory(q web.FIlterPaymentHistory) FilterPaymentHistory {
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

	return FilterPaymentHistory{
		Search: q.Search,
		UserId: q.UserId,
		Status: q.Status,
		Level:  q.Level,
		Pagination: Pagination{
			Page:  page,
			Limit: limit,
		},
	}
}

func (f *FilterPaymentHistory) QueryBuildFilterPayment() (filter, pagination string) {
	// Build base filter
	filter = BuildPaymentQuerySearch(f.Search)

	// Add level-based filter
	if f.Level == "2" {
		filter += fmt.Sprintf(" AND p.seller_id = '%s'", f.UserId)
	} else if f.Level == "1" {
		filter += fmt.Sprintf(" AND ph.user_id = '%s'", f.UserId)
	}

	// Add WHERE clause if filter exists
	if filter != "" {
		filter = filter[4:]
		filter = "WHERE " + filter
	}

	switch f.Status {
	case "1":
		filter += " AND ph.status = 'PENDING'"
	case "2":
		filter += " AND ph.status = 'SUCCEEDED'"
	case "3":
		filter += " AND ph.status = 'ACTIVATED'"
	case "4":
		filter += " AND ph.status = 'FAILED'"
	case "5":
		filter += " AND ph.status = 'EXPIRED'"
	}

	// Build pagination
	if f.Pagination.Page != "" && f.Pagination.Limit != "" {
		pagination = fmt.Sprintf("OFFSET %s LIMIT %s", f.Pagination.Page, f.Pagination.Limit)
	}

	return
}

func BuildPaymentQuerySearch(search string) (filter string) {
	fieldsAndValues := map[string]string{
		"p.name": search,
	}

	filter = helper.BuildManyColumnFilter(fieldsAndValues, filter)

	return
}
