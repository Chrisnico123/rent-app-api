package domain

import "time"

type TopUpRequest struct {
	Id        string
	UserId    string
	Amount    float64
	Type      string
	PaymentId string
	OrderID   string
	VAID      string
	ExpiredAt *time.Time
}

type TopUpResponse struct {
	Id        string
	Amount    float64
	Status    string
	Type      string
	PaymentId string
	OrderID   string
	VAID      string
	CreatedAt string
	UpdatedAt string
	ExpiredAt string
}

type TopUpByOrderIdResponse struct {
	Id        string
	Amount    float64
	Status    string
	Type      string
	PaymentId string
	OrderID   string
	VAID      string
	CreatedAt string
	UpdatedAt string
	ExpiredAt string
	UserId    string
	UserName  string
	UserEmail string
}
