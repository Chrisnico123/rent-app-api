package payment

import (
	"context"
	"fmt"
	"log"
	"math"
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

	CreateTopUP(ctx context.Context, request TopUpRequest, userId string) (web.PaymentTopUpByOrderIdResponse, error)
	GetListTopUP(ctx context.Context, userId string) ([]web.PaymentTopUpResponse, error)
	GetBalanceUser(ctx context.Context, userId string) (web.BalanceResponse, error)
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

// GetBalanceUser implements PaymentService.
func (s *paymentService) GetBalanceUser(ctx context.Context, userId string) (web.BalanceResponse, error) {
	data, err := s.paymentRepository.GetBalance(ctx, userId)
	if err != nil {
		return web.BalanceResponse{}, web.ErrInternalServer(err.Error())
	}

	return web.BalanceResponse{
		Balance: data,
	}, nil
}

// CreateTopUP implements PaymentService.
func (s *paymentService) CreateTopUP(ctx context.Context, request TopUpRequest, userId string) (web.PaymentTopUpByOrderIdResponse, error) {
	err := request.Validate()
	if err != nil {
		return web.PaymentTopUpByOrderIdResponse{}, web.ErrBadRequest(err.Error())
	}

	// Determine bank code based on TypeVA
	bankCode, err := helper.GetBankCode(request.TypeVA)
	if err != nil {
		return web.PaymentTopUpByOrderIdResponse{}, err
	}

	paymentType := helper.GetPaymentType(request.TypeVA)

	req := xendit.PaymentRequest{
		Amount:      request.Amount,
		BankCode:    bankCode,
		PaymentDesc: "TOPUP",
		OrderId:     fmt.Sprintf("Order-%d", time.Now().Unix()),
		ExpiredAt:   time.Now().Add(1 * time.Hour).UTC(),
	}

	resp, err := xendit.CreateTopUpVA(req)
	if err != nil {
		return web.PaymentTopUpByOrderIdResponse{}, fmt.Errorf("failed to create payment: %w", err)
	}

	topUpReq := domain.TopUpRequest{
		Id:        helper.GenerateId(),
		UserId:    userId,
		Amount:    req.Amount,
		Type:      paymentType,
		PaymentId: resp.PaymentMethod.Id,
		OrderID:   *resp.PaymentMethod.ReferenceId,
		VAID:      *resp.PaymentMethod.VirtualAccount.Get().ChannelProperties.VirtualAccountNumber,
		ExpiredAt: &req.ExpiredAt,
	}

	err = s.paymentRepository.CreateTopUp(ctx, topUpReq)
	if err != nil {
		return web.PaymentTopUpByOrderIdResponse{}, web.ErrInternalServer(err.Error())
	}

	data, err := s.paymentRepository.GetTopUpHistoryByOrderId(ctx, topUpReq.OrderID)
	if err != nil {
		return web.PaymentTopUpByOrderIdResponse{}, web.ErrInternalServer(err.Error())
	}

	result := web.PaymentTopUpByOrderIdResponse{
		Id:            data.Id,
		Amount:        data.Amount,
		Status:        data.Status,
		PaymentMethod: data.Type,
		PaymentId:     data.PaymentId,
		OrderId:       data.OrderID,
		VAID:          data.VAID,
		CreatedAt:     data.CreatedAt,
		UpdatedAt:     data.UpdatedAt,
		ExpiredAt:     data.ExpiredAt,
		Username:      data.UserName,
		Email:         data.UserEmail,
	}

	return result, nil
}

func (s *paymentService) GetListTopUP(ctx context.Context, userId string) ([]web.PaymentTopUpResponse, error) {
	// Get data from repository
	data, err := s.paymentRepository.GetTopUpHistory(ctx, userId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []web.PaymentTopUpResponse{}, nil
		}
		return nil, web.ErrInternalServer(err.Error())
	}

	// Map repository data to response struct
	responses := make([]web.PaymentTopUpResponse, 0, len(data))
	for _, payment := range data {
		response := web.PaymentTopUpResponse{
			Id:            payment.Id,
			Amount:        payment.Amount,
			Status:        payment.Status,
			PaymentMethod: payment.Type,
			PaymentId:     payment.PaymentId,
			VAID:          payment.VAID,
			OrderId:       payment.OrderID,
			CreatedAt:     payment.CreatedAt,
			UpdatedAt:     payment.UpdatedAt,
			ExpiredAt:     payment.ExpiredAt,
		}
		responses = append(responses, response)
	}

	return responses, nil
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
			OrderId:      v.OrderId,
			ProductName:  v.ProductName,
			CustomerName: v.CustomerName,
			Email:        v.Email,
			StartDate:    v.StartDate,
			EndDate:      v.EndDate,
			Price:        v.Price,
			Img:          v.Img[0],
			Status:       v.Status,
			Method:       v.Method,
			CreatedDate:  v.CreatedDate,
			ExpiredDate:  v.ExpiredDate,
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
		ID:           paymentHistory.ID,
		UserID:       paymentHistory.UserID,
		Username:     paymentHistory.Username,
		Email:        paymentHistory.Email,
		PaymentId:    paymentHistory.PaymentId,
		StartDate:    paymentHistory.StartDate,
		EndDate:      paymentHistory.EndDate,
		Method:       paymentHistory.Method,
		BookingDays:  int8(paymentHistory.BookingDays),
		TotalPayment: product.Price * float64(paymentHistory.BookingDays),
		OrderID:      paymentHistory.OrderID,
		VAID:         paymentHistory.VAID,
		Status:       paymentHistory.Status,
		CreatedAt:    paymentHistory.CreatedAt,
		UpdatedAt:    paymentHistory.UpdatedAt,
		ExpiredAt:    paymentHistory.ExpiredAt,
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

	// Process different events
	switch callbackPayload.Event {
	case "payment_method.expired":
		if callbackPayload.Data.VirtualAccount.ChannelProperties.CustomerName == "BOOK" {
			err := s.paymentRepository.UpdatePaymentStatus(ctx, callbackPayload.Data.ReferenceID, "EXPIRED")
			if err != nil {
				log.Println(err)
			}
		} else {
			err := s.paymentRepository.UpdateStatus(ctx, callbackPayload.Data.ReferenceID, "EXPIRED")
			if err != nil {
				log.Println(err)
			}
		}
		log.Printf("Payment expired: %s", callbackPayload.Data.ReferenceID)
	case "payment_method.activated":
		if callbackPayload.Data.VirtualAccount.ChannelProperties.CustomerName == "BOOK" {
			err := s.paymentRepository.UpdatePaymentStatus(ctx, callbackPayload.Data.ReferenceID, "ACTIVATED")
			if err != nil {
				log.Println(err)
			}
		} else {
			err := s.paymentRepository.UpdateStatus(ctx, callbackPayload.Data.ReferenceID, "ACTIVATED")
			if err != nil {
				log.Println(err)
			}
		}

		log.Printf("Payment activated: %s", callbackPayload.Data.ReferenceID)
	case "payment_method.failed":
		if callbackPayload.Data.VirtualAccount.ChannelProperties.CustomerName == "BOOK" {
			err := s.paymentRepository.UpdatePaymentStatus(ctx, callbackPayload.Data.ReferenceID, "FAILED")
			if err != nil {
				log.Println(err)
			}
		} else {
			err := s.paymentRepository.UpdateStatus(ctx, callbackPayload.Data.ReferenceID, "FAILED")
			if err != nil {
				log.Println(err)
			}
		}
		log.Printf("Payment failed: %s", callbackPayload.Data.ReferenceID)
	case "payment.succeeded":
		if callbackPayload.Data.PaymentMethod.VirtualAccount.ChannelProperties.CustomerName == "BOOK" {
			err := s.paymentRepository.AfterPaymentHandler(ctx, callbackPayload.Data.PaymentMethod.ReferenceID)
			if err != nil {
				return err
			}

		} else {
			err := s.paymentRepository.AfterPaymentTopUpHandler(ctx, callbackPayload.Data.PaymentMethod.ReferenceID)
			if err != nil {
				log.Println(err)
			}
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

	startDate, err := time.Parse("2006-01-02", req.StartdDate)
	if err != nil {
		return web.PaymentHistory{}, web.ErrBadRequest("invalid start date format, expected YYYY-MM-DD")
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return web.PaymentHistory{}, web.ErrBadRequest("invalid end date format, expected YYYY-MM-DD")
	}

	// Validate date range
	if endDate.Before(startDate) {
		return web.PaymentHistory{}, web.ErrBadRequest("end date cannot be before start date")
	}

	// Calculate duration and total price
	duration := endDate.Sub(startDate)
	days := int(math.Ceil(duration.Hours() / 24))
	if days < 1 {
		days = 1 // Minimum 1 day
	}
	totalPrice := product.Price * float64(days)

	var orderReq domain.OrderUser

	// Check If using balance
	if req.TypeVA == 4 {
		paymentReq := domain.PaymentHistory{
			ID:        helper.GenerateId(),
			UserID:    userId,
			PaymentId: "",
			OrderID:   helper.GenerateOrderId(),
			VAID:      "",
			Status:    "SUCCEEDED",
			ExpiredAt: nil,
		}

		orderReq = domain.OrderUser{
			ID:        helper.GenerateId(),
			UserID:    userId,
			PaymentId: "",
			StartDate: req.StartdDate,
			EndDate:   req.EndDate,
			Desc:      req.Desc,
			OrderID:   paymentReq.OrderID,
			ProductID: &req.ProductId,
			Type:      "Balance",
			Price:     totalPrice,
		}

		balance, err := s.paymentRepository.GetBalance(ctx, userId)
		if err != nil {
			return web.PaymentHistory{}, web.ErrInternalServer(err.Error())
		}

		if balance < totalPrice {
			return web.PaymentHistory{}, web.ErrBadRequest("balance not enough, please top up your balance")
		}

		err = s.paymentRepository.CreateOrderBalance(ctx, orderReq, paymentReq, reqImg, reqIvs, totalPrice)
		if err != nil {
			return web.PaymentHistory{}, web.ErrInternalServer(err.Error())
		}

		err = s.paymentRepository.AfterPaymentHandler(ctx, orderReq.OrderID)
		if err != nil {
			return web.PaymentHistory{}, web.ErrInternalServer(err.Error())
		}
	} else {
		// Determine bank code based on TypeVA
		bankCode, err := helper.GetBankCode(req.TypeVA)
		if err != nil {
			return web.PaymentHistory{}, err
		}

		paymentType := helper.GetPaymentType(req.TypeVA)

		request := xendit.PaymentRequest{
			Amount:      totalPrice,
			BankCode:    bankCode,
			PaymentDesc: "BOOK",
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

		orderReq = domain.OrderUser{
			ID:        helper.GenerateId(),
			UserID:    userId,
			PaymentId: resp.PaymentMethod.Id,
			StartDate: req.StartdDate,
			EndDate:   req.EndDate,
			Desc:      req.Desc,
			OrderID:   *resp.PaymentMethod.ReferenceId,
			ProductID: &req.ProductId,
			Type:      paymentType,
			Price:     totalPrice,
		}

		err = s.paymentRepository.CreateOrder(ctx, orderReq, paymentReq, reqImg, reqIvs)
		if err != nil {
			return web.PaymentHistory{}, web.ErrInternalServer(err.Error())
		}
	}

	data, err := s.GetPaymentByOrderId(ctx, orderReq.OrderID)
	if err != nil {
		return web.PaymentHistory{}, web.ErrInternalServer(err.Error())
	}

	return data, nil
}
