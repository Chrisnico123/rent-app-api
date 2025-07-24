package payment

import (
	"context"
	"rent-application/apps/product"
	"rent-application/domain"
	"rent-application/internal/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentRepository interface {
	CreateOrder(ctx context.Context, order domain.OrderUser, payment domain.PaymentHistory) error
	AfterPaymentHandler(ctx context.Context, orderId string) error

	CreatePayment(ctx context.Context, payment domain.PaymentHistory) error
	GetPaymentByID(ctx context.Context, id string) (domain.PaymentHistory, error)
	GetOrdersByUserID(ctx context.Context, userID string, filter domain.FilterSearchPagination) ([]domain.OrderUser, int, error)
	GetPaymentByOrderID(ctx context.Context, orderId string) (domain.PaymentResponse, error)
	GetPaymentsByOrderID(ctx context.Context, orderID string) (domain.PaymentHistory, error)
	UpdatePaymentStatus(ctx context.Context, paymentID string, status string) error
}

type paymentRepository struct {
	db           database.Store
	paymentQuery PaymentQuery
	productQuery product.ProductQuery
}

func NewPaymentRepository(db database.Store, paymentQuery PaymentQuery, productQuery product.ProductQuery) PaymentRepository {
	return &paymentRepository{
		db:           db,
		paymentQuery: paymentQuery,
		productQuery: productQuery,
	}
}

// GetPaymentByOrderID implements PaymentRepository.
func (r *paymentRepository) GetPaymentByOrderID(ctx context.Context, orderId string) (domain.PaymentResponse, error) {
	var payment domain.PaymentResponse
	err := r.db.WithoutTransaction(ctx, func(db *pgxpool.Pool) error {
		var err error
		payment, err = r.paymentQuery.GetPaymentByOrderID(ctx, db, orderId)
		return err
	})
	return payment, err
}

// AfterPaymentHandler implements PaymentRepository.
func (r *paymentRepository) AfterPaymentHandler(ctx context.Context, orderId string) error {
	var data domain.OrderUser
	var err error
	return r.db.WithTransaction(ctx, func(tx pgx.Tx) error {
		if data, err = r.paymentQuery.GetOrderByID(ctx, tx, orderId); err != nil {
			return err
		}

		if err = r.productQuery.UpdateAvailableProduct(ctx, tx, *data.ProductID, false); err != nil {
			return err
		}

		if err := r.paymentQuery.UpdatePaymentStatus(ctx, tx, orderId, "SUCCEEDED"); err != nil {
			return err
		}

		return nil
	})
}

func (r *paymentRepository) CreateOrder(ctx context.Context, order domain.OrderUser, payment domain.PaymentHistory) error {
	return r.db.WithTransaction(ctx, func(tx pgx.Tx) error {
		if err := r.paymentQuery.CreatePayment(ctx, tx, payment); err != nil {
			return err
		}

		if err := r.paymentQuery.CreateBooking(ctx, tx, order); err != nil {
			return err
		}

		return nil
	})
}

func (r *paymentRepository) CreatePayment(ctx context.Context, payment domain.PaymentHistory) error {
	return r.db.WithTransaction(ctx, func(tx pgx.Tx) error {
		return r.paymentQuery.CreatePayment(ctx, tx, payment)
	})
}

func (r *paymentRepository) GetPaymentByID(ctx context.Context, id string) (domain.PaymentHistory, error) {
	var payment domain.PaymentHistory
	err := r.db.WithoutTransaction(ctx, func(db *pgxpool.Pool) error {
		var err error
		// Assuming you'll implement GetPaymentByID in PaymentQuery
		payment, err = r.paymentQuery.GetPaymentByID(ctx, db, id)
		return err
	})
	return payment, err
}

func (r *paymentRepository) GetOrdersByUserID(ctx context.Context, userID string, filter domain.FilterSearchPagination) ([]domain.OrderUser, int, error) {
	var orders []domain.OrderUser
	var total int

	err := r.db.WithoutTransaction(ctx, func(db *pgxpool.Pool) error {
		var err error
		// Assuming you'll implement these in PaymentQuery
		orders, err = r.paymentQuery.GetOrdersByUserID(ctx, db, userID, filter)
		if err != nil {
			return err
		}

		total, err = r.paymentQuery.CountOrdersByUserID(ctx, db, userID, filter)
		return err
	})

	return orders, total, err
}

func (r *paymentRepository) GetPaymentsByOrderID(ctx context.Context, orderID string) (domain.PaymentHistory, error) {
	var payments domain.PaymentHistory
	err := r.db.WithoutTransaction(ctx, func(db *pgxpool.Pool) error {
		var err error
		// Assuming you'll implement this in PaymentQuery
		payments, err = r.paymentQuery.GetPaymentsByOrderID(ctx, db, orderID)
		return err
	})
	return payments, err
}

func (r *paymentRepository) UpdatePaymentStatus(ctx context.Context, paymentID string, status string) error {
	return r.db.WithTransaction(ctx, func(tx pgx.Tx) error {
		// Assuming you'll implement this in PaymentQuery
		return r.paymentQuery.UpdatePaymentStatus(ctx, tx, paymentID, status)
	})
}
