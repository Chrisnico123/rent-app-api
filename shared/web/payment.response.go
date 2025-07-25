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
	OrderID       string        `json:"order_id"`
	ProductDetail ProductDetail `json:"product_detail"`
	VAID          string        `json:"virtual_account_id"`
	Status        string        `json:"status"`
	CreatedAt     string        `json:"created_at"`
	UpdatedAt     string        `json:"updated_at"`
	ExpiredAt     string        `json:"expired_at"`
}
