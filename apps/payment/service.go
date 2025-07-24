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

	CreateBooking(ctx context.Context, userId string, req web.BookRequest) (web.PaymentHistory, error)
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
		// Handle expired payment
		err := s.paymentRepository.UpdatePaymentStatus(ctx, callbackPayload.Data.ReferenceID, callbackPayload.Data.Status)
		if err != nil {
			log.Println(err)
		}

		log.Printf("Payment expired: %s", callbackPayload.Data.ReferenceID)
	case "payment_method.activated":
		log.Printf("Payment activated: %s", callbackPayload.Data.ReferenceID)
	case "payment_method.failed":
		// Handle failed payment
		log.Printf("Payment failed: %s", callbackPayload.Data.ReferenceID)
	case "payment.succeeded":
		// Handle activated payment
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

func (s *paymentService) CreateBooking(ctx context.Context, userId string, req web.BookRequest) (web.PaymentHistory, error) {
	product, err := s.productRepository.GetProductByID(ctx, req.ProductId)
	if err != nil {
		return web.PaymentHistory{}, fmt.Errorf("failed to get product: %w", err)
	}
	if product.ID == "" {
		return web.PaymentHistory{}, web.ErrNotFound("product not found")
	}

	if !product.Available {
		return web.PaymentHistory{}, web.ErrBadRequest("product not available, please check another product")
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
		OrderID:   *resp.PaymentMethod.ReferenceId,
		ProductID: &req.ProductId,
		Type:      paymentType,
		Price:     product.Price,
	}

	err = s.paymentRepository.CreateOrder(ctx, orderReq, paymentReq)
	if err != nil {
		return web.PaymentHistory{}, web.ErrInternalServer(err.Error())
	}

	paymentHistory, err := s.paymentRepository.GetPaymentByOrderID(ctx, paymentReq.OrderID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return web.PaymentHistory{}, web.ErrNotFound("payment not found")
		}
		return web.PaymentHistory{}, web.ErrInternalServer(err.Error())
	}

	res := web.PaymentHistory{
		ID:        paymentHistory.ID,
		UserID:    userId,
		PaymentId: paymentHistory.PaymentId,
		OrderID:   paymentHistory.OrderID,
		Status:    paymentHistory.Status,
		ProductDetail: web.ProductDetail{
			Name:         product.Name,
			Price:        product.Price,
			Description:  product.Description,
			CategoryName: product.CategoryName,
			Img:          product.Img[0],
		},
		CreatedAt: paymentHistory.CreatedAt,
		UpdatedAt: paymentHistory.UpdatedAt,
		ExpiredAt: paymentHistory.ExpiredAt,
	}

	return res, nil
}
