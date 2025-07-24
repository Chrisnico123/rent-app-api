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
			ReferenceID string `json:"reference_id"`
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
		// Card               interface{} `json:"card"`
		// Ewallet            interface{} `json:"ewallet"`
		// QrCode             interface{} `json:"qr_code"`
		// Metadata           interface{} `json:"metadata"`
		// CustomerID         interface{} `json:"customer_id"`
		// Description        interface{} `json:"description"`
		// Reusability        string      `json:"reusability"`
		// DirectDebit        interface{} `json:"direct_debit"`
		// FailureCode        interface{} `json:"failure_code"`
		// OverTheCounter     interface{} `json:"over_the_counter"`
		// DirectBankTransfer interface{} `json:"direct_bank_transfer"`
		// BillingInformation struct {
		// 	City          interface{} `json:"city"`
		// 	Country       string      `json:"country"`
		// 	PostalCode    interface{} `json:"postal_code"`
		// 	StreetLine1   interface{} `json:"street_line1"`
		// 	StreetLine2   interface{} `json:"street_line2"`
		// 	ProvinceState interface{} `json:"province_state"`
		// } `json:"billing_information"`
		// Actions []interface{} `json:"actions"`
	} `json:"data"`
}
