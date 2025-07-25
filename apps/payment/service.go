package payment

import (
	"context"
	"fmt"
	"log"
	"rent-application/apps/auth"
	"rent-application/apps/product"
	"rent-application/domain"
	"rent-application/internal/payment/xendit"
	"rent-application/shared/helper"
	"rent-application/shared/web"
	"time"

	"github.com/jackc/pgx/v5"
)

type PaymentService interface {
	GetPaymentByOrderId(ctx context.Context, orderId string) (web.PaymentHistory, error)
	GetPaymentListPayment(ctx context.Context, filter web.FIlterPaymentHistory) ([]web.PaymentHistoryResponse, int, error)

	CreateBooking(ctx context.Context, userId string, req BookRequest) (web.PaymentHistory, error)
	AfterPaymentHandler(ctx context.Context, request web.XenditVACallback) error
}

type paymentService struct {
	paymentRepository PaymentRepository
	productRepository product.ProductRepository
	authRepository    auth.AuthRepository
}

func NewPaymentService(paymentRepository PaymentRepository, productRepository product.ProductRepository, authRepository auth.AuthRepository) PaymentService {
	return &paymentService{
		paymentRepository: paymentRepository,
		productRepository: productRepository,
		authRepository:    authRepository,
	}
}

// GetPaymentListPayment implements PaymentService.
func (s *paymentService) GetPaymentListPayment(ctx context.Context, filter web.FIlterPaymentHistory) ([]web.PaymentHistoryResponse, int, error) {
	// Convert web filter to domain filter
	domainFilter := domain.ToDomainFilterPaymentHistory(filter)

	// Get data from repository
	datas, count, err := s.paymentRepository.GetListPaymentHistory(ctx, domainFilter)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []web.PaymentHistoryResponse{}, 0, nil
		}
		return []web.PaymentHistoryResponse{}, 0, web.ErrInternalServer(err.Error())
	}

	// Convert domain models to web responses
	list := make([]web.PaymentHistoryResponse, 0, len(datas))
	for _, v := range datas {
		data := web.PaymentHistoryResponse{
			OrderId:     v.OrderId,
			ProductName: v.ProductName,
			Price:       v.Price,
			Img:         v.Img[0],
			Status:      v.Status,
			Method:      v.Method,
			CreatedDate: v.CreatedDate,
			ExpiredDate: v.ExpiredDate,
		}
		list = append(list, data)
	}

	return list, count, nil
}

func (s *paymentService) GetPaymentByOrderId(ctx context.Context, orderId string) (web.PaymentHistory, error) {
	// Get payment data from repository
	paymentHistory, err := s.paymentRepository.GetPaymentByOrderID(ctx, orderId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return web.PaymentHistory{}, web.ErrNotFound("payment not found")
		}
		return web.PaymentHistory{}, web.ErrInternalServer(err.Error())
	}

	// Get product details
	product, err := s.productRepository.GetProductByID(ctx, paymentHistory.ProductId)
	if err != nil {
		return web.PaymentHistory{}, web.ErrInternalServer(fmt.Sprintf("failed to get product details: %v", err))
	}

	// Convert to response format
	response := web.PaymentHistory{
		ID:        paymentHistory.ID,
		UserID:    paymentHistory.UserID,
		Username:  paymentHistory.Username,
		Email:     paymentHistory.Email,
		PaymentId: paymentHistory.PaymentId,
		OrderID:   paymentHistory.OrderID,
		VAID:      paymentHistory.VAID,
		Status:    paymentHistory.Status,
		CreatedAt: paymentHistory.CreatedAt,
		UpdatedAt: paymentHistory.UpdatedAt,
		ExpiredAt: paymentHistory.ExpiredAt,
		ProductDetail: web.ProductDetail{
			Name:         product.Name,
			Price:        product.Price,
			Description:  product.Description,
			CategoryName: product.CategoryName,
			Img:          product.Img[0],
		},
	}

	return response, nil
}

// AfterPaymentHandler implements PaymentService.
func (s *paymentService) AfterPaymentHandler(ctx context.Context, callbackPayload web.XenditVACallback) error {
	// Log important information
	log.Printf("Callback Event: %s", callbackPayload.Event)
	log.Printf("Payment Method ID: %s", callbackPayload.Data.ID)
	log.Printf("Reference ID: %s", callbackPayload.Data.ReferenceID)
	log.Printf("Status: %s", callbackPayload.Data.Status)

	if callbackPayload.Data.VirtualAccount != nil {
		log.Printf("Virtual Account Details:")
		log.Printf("- Amount: %.2f", callbackPayload.Data.VirtualAccount.Amount)
		log.Printf("- Channel Code: %s", callbackPayload.Data.VirtualAccount.ChannelCode)
		log.Printf("- VA Number: %s", callbackPayload.Data.VirtualAccount.ChannelProperties.VirtualAccountNumber)
	}

	// Process different events
	switch callbackPayload.Event {
	case "payment_method.expired":
		err := s.paymentRepository.UpdatePaymentStatus(ctx, callbackPayload.Data.ReferenceID, "EXPIRED")
		if err != nil {
			log.Println(err)
		}
		log.Printf("Payment expired: %s", callbackPayload.Data.ReferenceID)
	case "payment_method.activated":
		err := s.paymentRepository.UpdatePaymentStatus(ctx, callbackPayload.Data.ReferenceID, "ACTIVATED")
		if err != nil {
			log.Println(err)
		}
		log.Printf("Payment activated: %s", callbackPayload.Data.ReferenceID)
	case "payment_method.failed":
		err := s.paymentRepository.UpdatePaymentStatus(ctx, callbackPayload.Data.ReferenceID, "FAILED")
		if err != nil {
			log.Println(err)
		}
		log.Printf("Payment failed: %s", callbackPayload.Data.ReferenceID)
	case "payment.succeeded":
		err := s.paymentRepository.AfterPaymentHandler(ctx, callbackPayload.Data.PaymentMethod.ReferenceID)
		if err != nil {
			return err
		}
		log.Printf("Payment succeeded: %s", callbackPayload.Data.PaymentMethod.ReferenceID)
	default:
		log.Printf("Unhandled event type: %s", callbackPayload.Event)
	}

	return nil
}

func (s *paymentService) CreateBooking(ctx context.Context, userId string, req BookRequest) (web.PaymentHistory, error) {
	err := req.Validate()
	if err != nil {
		return web.PaymentHistory{}, web.ErrBadRequest(err.Error())
	}
	product, err := s.productRepository.GetProductByID(ctx, req.ProductId)

	if err != nil {
		return web.PaymentHistory{}, web.ErrInternalServer(err.Error())
	}
	if product.ID == "" {
		return web.PaymentHistory{}, web.ErrNotFound("product not found")
	}

	if !product.Available {
		return web.PaymentHistory{}, web.ErrBadRequest("product not available, please check another product")
	}

	var reqImg domain.EncryptedImg
	var reqIvs domain.EncryptedIv

	if req.File != "" {
		reqImg.Id = helper.GenerateId()
		reqImg.UserId = userId

		reqIvs.Id = helper.GenerateId()
		reqIvs.EncryId = reqImg.Id

		// Encrypted Data
		reqImg.EncryUrl, reqIvs.Ivs, err = helper.Encrypt(req.File)
		if err != nil {
			return web.PaymentHistory{}, err
		}
	}

	// Determine bank code based on TypeVA
	bankCode, err := helper.GetBankCode(req.TypeVA)
	if err != nil {
		return web.PaymentHistory{}, err
	}

	paymentType := helper.GetPaymentType(req.TypeVA)

	request := xendit.PaymentRequest{
		Amount:      product.Price,
		BankCode:    bankCode,
		PaymentDesc: fmt.Sprintf("Payment for product %s", product.Name),
		OrderId:     fmt.Sprintf("Order-%d", time.Now().Unix()),
		ExpiredAt:   time.Now().Add(1 * time.Hour).UTC(),
	}

	resp, err := xendit.CreatePaymentVA(request)
	if err != nil {
		return web.PaymentHistory{}, fmt.Errorf("failed to create payment: %w", err)
	}

	paymentReq := domain.PaymentHistory{
		ID:        helper.GenerateId(),
		UserID:    userId,
		PaymentId: resp.PaymentMethod.Id,
		OrderID:   *resp.PaymentMethod.ReferenceId,
		VAID:      *resp.PaymentMethod.VirtualAccount.Get().ChannelProperties.VirtualAccountNumber,
		Status:    "PENDING",
		ExpiredAt: &request.ExpiredAt,
	}

	orderReq := domain.OrderUser{
		ID:        helper.GenerateId(),
		UserID:    userId,
		PaymentId: resp.PaymentMethod.Id,
		StartDate: req.StartdDate,
		EndDate:   req.EndDate,
		Desc:      req.Desc,
		OrderID:   *resp.PaymentMethod.ReferenceId,
		ProductID: &req.ProductId,
		Type:      paymentType,
		Price:     product.Price,
	}

	err = s.paymentRepository.CreateOrder(ctx, orderReq, paymentReq, reqImg, reqIvs)
	if err != nil {
		return web.PaymentHistory{}, web.ErrInternalServer(err.Error())
	}

	data, err := s.GetPaymentByOrderId(ctx, orderReq.OrderID)
	if err != nil {
		return web.PaymentHistory{}, web.ErrInternalServer(err.Error())
	}

	return data, nil
}
