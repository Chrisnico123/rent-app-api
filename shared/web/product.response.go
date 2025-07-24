package web

type ProductResponse struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Description  interface{} `json:"description"`
	Img          interface{} `json:"img"`
	CategoryID   string      `json:"category_id"`
	CategoryName string      `json:"category_name"`
	Price        float64     `json:"price"`
	Available    bool        `json:"available"`
	SellerID     string      `json:"seller_id"`
	CreatedAt    string      `json:"created_at"`
	UpdatedAt    string      `json:"updated_at"`
}

type ProductResponseList struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Img       string  `json:"img"`
	Price     float64 `json:"price"`
	Available bool    `json:"available"`
	Category  string  `json:"category"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}
