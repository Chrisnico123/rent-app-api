package domain

import "time"

type OrderUser struct {
	ID        string
	UserID    string
	PaymentId string
	OrderID   string
	ProductID *string
	Type      string
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
	ID        string
	UserID    string
	PaymentId string
	OrderID   string
	ProductId string
	VAID      string
	Status    string
	CreatedAt string
	UpdatedAt string
	ExpiredAt string
}
