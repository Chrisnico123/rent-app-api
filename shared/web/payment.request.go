package web

type PaymentRequest struct {
	Amount uint16 `query:"amount"`
	TypeVA uint8  `query:"type"`
}

type FIlterPaymentHistory struct {
	Search string `query:"search"`
	Limit  string `query:"limit"`
	Page   string `query:"page"`
	Status string `query:"status"`
	Level  string
	UserId string
}

type PaymentHistoryResponse struct {
	OrderId      string `json:"order_id"`
	CustomerName string `json:"customer_name"`
	Email        string `json:"email"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	ProductName  string `json:"name"`
	Price        string `json:"price"`
	Img          string `json:"img"`
	Status       string `json:"status"`
	Method       string `json:"method"`
	CreatedDate  string `json:"book_date"`
	ExpiredDate  string `json:"expired_date"`
}
