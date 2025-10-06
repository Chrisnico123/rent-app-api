package web

type ProductDetail struct {
	Name         string   `json:"name"`
	Price        float64  `json:"price"`
	Description  []string `json:"description"`
	CategoryName string   `json:"category"`
	Img          string   `json:"img"`
}

type PaymentHistory struct {
	ID            string        `json:"id"`
	UserID        string        `json:"user_id"`
	Username      string        `json:"username"`
	Email         string        `json:"email"`
	PaymentId     string        `json:"payment_id"`
	Method        string        `json:"method"`
	OrderID       string        `json:"order_id"`
	StartDate     string        `json:"start_date"`
	EndDate       string        `json:"end_date"`
	BookingDays   int8          `json:"booking_days"`
	TotalPayment  float64       `json:"total_payment"`
	ProductDetail ProductDetail `json:"product_detail"`
	VAID          string        `json:"virtual_account_id"`
	Status        string        `json:"status"`
	CreatedAt     string        `json:"created_at"`
	UpdatedAt     string        `json:"updated_at"`
	ExpiredAt     string        `json:"expired_at"`
}

type PaymentTopUpResponse struct {
	Id            string  `json:"id"`
	Amount        float64 `json:"amount"`
	Status        string  `json:"status"`
	PaymentMethod string  `json:"payment_method"`
	PaymentId     string  `json:"payment_id"`
	VAID          string  `json:"virtual_account_id"`
	OrderId       string  `json:"order_id"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	ExpiredAt     string  `json:"expired_at"`
}

type PaymentTopUpByOrderIdResponse struct {
	Id            string  `json:"id"`
	Amount        float64 `json:"amount"`
	Status        string  `json:"status"`
	PaymentMethod string  `json:"payment_method"`
	PaymentId     string  `json:"payment_id"`
	VAID          string  `json:"virtual_account_id"`
	OrderId       string  `json:"order_id"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	ExpiredAt     string  `json:"expired_at"`
	Username      string  `json:"username"`
	Email         string  `json:"email"`
}

type BalanceResponse struct {
	Balance float64 `json:"balance"`
}
