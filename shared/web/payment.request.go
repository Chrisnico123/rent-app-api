package web

type PaymentRequest struct {
	Amount uint16 `query:"amount"`
	TypeVA uint8  `query:"type"`
}

type BookRequest struct {
	ProductId string `query:"product_id"`
	TypeVA    uint8  `query:"type"`
}
