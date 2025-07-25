package payment

import (
	"errors"
	"time"
	"unicode/utf8"
)

type BookRequest struct {
	ProductId  string `json:"product_id"`
	StartdDate string `json:"start_date"`
	File       string `json:"file"`
	EndDate    string `json:"end_date"`
	Desc       string `json:"desc"`
	TypeVA     uint8  `json:"type"`
}

func (r *BookRequest) Validate() error {
	// Validate ProductId
	if r.ProductId == "" {
		return errors.New("product_id is required")
	}
	if utf8.RuneCountInString(r.ProductId) > 15 {
		return errors.New("product_id is too long (max 15 characters)")
	}

	// Validate StartDate
	if r.StartdDate == "" {
		return errors.New("start_date is required")
	}
	startDate, err := time.Parse("2006-01-02", r.StartdDate)
	if err != nil {
		return errors.New("start_date must be in YYYY-MM-DD format")
	}

	// Validate EndDate
	if r.EndDate == "" {
		return errors.New("end_date is required")
	}
	endDate, err := time.Parse("2006-01-02", r.EndDate)
	if err != nil {
		return errors.New("end_date must be in YYYY-MM-DD format")
	}

	// Validate date range
	if endDate.Before(startDate) {
		return errors.New("end_date cannot be before start_date")
	}

	// Validate Description
	if utf8.RuneCountInString(r.Desc) > 500 {
		return errors.New("description is too long (max 500 characters)")
	}

	// Validate TypeVA
	if r.TypeVA < 1 || r.TypeVA > 3 {
		return errors.New("type must be between 1 and 3")
	}

	return nil
}
