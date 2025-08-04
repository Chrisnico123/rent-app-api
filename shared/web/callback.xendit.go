package web

type XenditVACallback struct {
	ID         string `json:"id"`
	Event      string `json:"event"`
	Created    string `json:"created"`
	BusinessID string `json:"business_id"`
	Data       struct {
		ID            string `json:"id"`
		Type          string `json:"type"`
		Status        string `json:"status"`
		Country       string `json:"country"`
		Created       string `json:"created"`
		Updated       string `json:"updated"`
		ReferenceID   string `json:"reference_id"`
		PaymentMethod *struct {
			ReferenceID    string `json:"reference_id"`
			VirtualAccount *struct {
				Amount            float64 `json:"amount"`
				Currency          string  `json:"currency"`
				ChannelCode       string  `json:"channel_code"`
				ChannelProperties struct {
					CustomerName         string `json:"customer_name"`
					ExpiresAt            string `json:"expires_at"`
					VirtualAccountNumber string `json:"virtual_account_number"`
				} `json:"channel_properties"`
			} `json:"virtual_account"`
		} `json:"payment_method"`
		VirtualAccount *struct {
			Amount            float64 `json:"amount"`
			Currency          string  `json:"currency"`
			ChannelCode       string  `json:"channel_code"`
			ChannelProperties struct {
				CustomerName         string `json:"customer_name"`
				ExpiresAt            string `json:"expires_at"`
				VirtualAccountNumber string `json:"virtual_account_number"`
			} `json:"channel_properties"`
		} `json:"virtual_account"`
	} `json:"data"`
}
