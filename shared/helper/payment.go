package helper

import (
	"rent-application/internal/payment/xendit"
	"rent-application/shared/constants"
	"rent-application/shared/web"
)

// Helper function to get bank code from TypeVA
func GetBankCode(typeVA uint8) (string, error) {
	switch typeVA {
	case constants.PaymentTypeMandiri:
		return xendit.BankMandiri, nil
	case constants.PaymentTypeBCA:
		return xendit.BankBCA, nil
	case constants.PaymentTypeBRI:
		return xendit.BankBRI, nil
	default:
		return "", web.ErrBadRequest("invalid payment type")
	}
}

// Helper function to get payment type string
func GetPaymentType(typeVA uint8) string {
	switch typeVA {
	case constants.PaymentTypeMandiri:
		return "VA-MANDIRI"
	case constants.PaymentTypeBCA:
		return "VA-BCA"
	case constants.PaymentTypeBRI:
		return "VA-BRI"
	default:
		return "BOOK-UNKNOWN"
	}
}
