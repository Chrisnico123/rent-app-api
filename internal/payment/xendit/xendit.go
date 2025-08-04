package xendit

import (
	"context"
	"fmt"
	"rent-application/configs"
	"rent-application/shared/web"
	"time"

	"github.com/xendit/xendit-go/v7"
	"github.com/xendit/xendit-go/v7/payment_request"
)

const (
	BankMandiri = "MANDIRI"
	BankBCA     = "BCA"
	BankBRI     = "BRI"
)

func float64Ptr(f float64) *float64 { return &f }
func stringPtr(s string) *string    { return &s }

func InitXendit(apiKey string) *xendit.APIClient {
	return xendit.NewClient(apiKey)
}

type PaymentRequest struct {
	Amount       float64
	BankCode     string
	CustomerName string
	PaymentDesc  string
	OrderId      string
	ExpiredAt    time.Time
}

func CreatePaymentVA(request PaymentRequest) (*payment_request.PaymentRequest, error) {
	xenditClient := InitXendit(configs.Cfg.Xendit.ApiKey)
	var channelCode payment_request.VirtualAccountChannelCode

	switch request.BankCode {
	case "MANDIRI":
		channelCode = payment_request.VIRTUALACCOUNTCHANNELCODE_MANDIRI
	case "BCA":
		channelCode = payment_request.VIRTUALACCOUNTCHANNELCODE_BCA
	case "BRI":
		channelCode = payment_request.VIRTUALACCOUNTCHANNELCODE_BRI
	default:
		return nil, web.ErrBadRequest(fmt.Sprintf("unsupported bank code: %s", request.BankCode))
	}

	vaParams := payment_request.VirtualAccountParameters{
		ChannelCode: channelCode,
		ChannelProperties: payment_request.VirtualAccountChannelProperties{
			CustomerName: "BOOK",
			ExpiresAt:    &request.ExpiredAt,
		},
	}

	params := payment_request.PaymentRequestParameters{
		Amount:      float64Ptr(request.Amount),
		Currency:    payment_request.PAYMENTREQUESTCURRENCY_IDR,
		ReferenceId: stringPtr(fmt.Sprintf("BOOK-%d", int64(request.Amount))),
		Description: *payment_request.NewNullableString(&request.PaymentDesc),
		PaymentMethod: &payment_request.PaymentMethodParameters{
			Type:           payment_request.PAYMENTMETHODTYPE_VIRTUAL_ACCOUNT,
			VirtualAccount: *payment_request.NewNullableVirtualAccountParameters(&vaParams),
			Reusability:    payment_request.PAYMENTMETHODREUSABILITY_ONE_TIME_USE,
		},
	}

	fmt.Println("LOG - > ", params)

	resp, _, err := xenditClient.PaymentRequestApi.CreatePaymentRequest(context.Background()).
		PaymentRequestParameters(params).
		Execute()

	if err != nil {
		return nil, fmt.Errorf("xendit api error: %v", err)
	}

	return resp, nil
}

func CreateTopUpVA(request PaymentRequest) (*payment_request.PaymentRequest, error) {
	xenditClient := InitXendit(configs.Cfg.Xendit.ApiKey)
	var channelCode payment_request.VirtualAccountChannelCode

	switch request.BankCode {
	case "MANDIRI":
		channelCode = payment_request.VIRTUALACCOUNTCHANNELCODE_MANDIRI
	case "BCA":
		channelCode = payment_request.VIRTUALACCOUNTCHANNELCODE_BCA
	case "BRI":
		channelCode = payment_request.VIRTUALACCOUNTCHANNELCODE_BRI
	default:
		return nil, web.ErrBadRequest(fmt.Sprintf("unsupported bank code: %s", request.BankCode))
	}

	vaParams := payment_request.VirtualAccountParameters{
		ChannelCode: channelCode,
		ChannelProperties: payment_request.VirtualAccountChannelProperties{
			CustomerName: "TOPUP",
			ExpiresAt:    &request.ExpiredAt,
		},
	}

	var description *payment_request.NullableString
	if request.PaymentDesc != "" {
		description = payment_request.NewNullableString(&request.PaymentDesc)
	}

	params := payment_request.PaymentRequestParameters{
		Amount:      float64Ptr(request.Amount),
		Currency:    payment_request.PAYMENTREQUESTCURRENCY_IDR,
		ReferenceId: stringPtr(fmt.Sprintf("TOPUP-%d", int64(request.Amount))), // Ganti sesuai kebutuhan
		Description: *description,
		PaymentMethod: &payment_request.PaymentMethodParameters{
			Type:           payment_request.PAYMENTMETHODTYPE_VIRTUAL_ACCOUNT,
			VirtualAccount: *payment_request.NewNullableVirtualAccountParameters(&vaParams),
			Reusability:    payment_request.PAYMENTMETHODREUSABILITY_ONE_TIME_USE,
		},
	}

	resp, _, err := xenditClient.PaymentRequestApi.CreatePaymentRequest(context.Background()).
		PaymentRequestParameters(params).
		Execute()

	if err != nil {
		return nil, fmt.Errorf("xendit api error: %v", err)
	}

	return resp, nil
}
