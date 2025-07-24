package payment

import (
	"context"
	"fmt"
	"rent-application/domain"
	"rent-application/shared/helper"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentQuery interface {
	// Order operations
	CreateBooking(ctx context.Context, tx pgx.Tx, req domain.OrderUser) error
	GetOrderByID(ctx context.Context, tx pgx.Tx, id string) (domain.OrderUser, error)
	GetOrdersByUserID(ctx context.Context, db *pgxpool.Pool, userID string, filter domain.FilterSearchPagination) ([]domain.OrderUser, error)
	CountOrdersByUserID(ctx context.Context, db *pgxpool.Pool, userID string, filter domain.FilterSearchPagination) (int, error)

	// Payment operations
	CreatePayment(ctx context.Context, tx pgx.Tx, req domain.PaymentHistory) error
	GetPaymentByID(ctx context.Context, db *pgxpool.Pool, id string) (domain.PaymentHistory, error)
	GetPaymentsByOrderID(ctx context.Context, db *pgxpool.Pool, orderID string) (domain.PaymentHistory, error)
	UpdatePaymentStatus(ctx context.Context, tx pgx.Tx, paymentID string, status string) error
	GetPaymentByOrderID(ctx context.Context, db *pgxpool.Pool, id string) (domain.PaymentResponse, error)
}

type PaymentQueryImpl struct{}

func NewPaymentQuery() PaymentQuery {
	return &PaymentQueryImpl{}
}

// Order operations
func (q *PaymentQueryImpl) CreateBooking(ctx context.Context, tx pgx.Tx, req domain.OrderUser) error {
	query := `
		INSERT INTO order_user (
			id, 
			user_id, 
			order_id, 
			product_id, 
			type, 
			price
		) VALUES (
			$1, $2, $3, $4, $5, $6
		)`

	_, err := tx.Exec(ctx, query,
		req.ID,
		req.UserID,
		req.OrderID,
		req.ProductID,
		req.Type,
		req.Price,
	)
	return err
}

func (q *PaymentQueryImpl) GetOrderByID(ctx context.Context, tx pgx.Tx, id string) (domain.OrderUser, error) {
	query := `
		SELECT 
			id, user_id, order_id, product_id, type, price, created_at, updated_at
		FROM 
			order_user
		WHERE 
			order_id = $1`

	var order domain.OrderUser
	err := tx.QueryRow(ctx, query, id).Scan(
		&order.ID,
		&order.UserID,
		&order.OrderID,
		&order.ProductID,
		&order.Type,
		&order.Price,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	return order, err
}

func (q *PaymentQueryImpl) GetOrdersByUserID(ctx context.Context, db *pgxpool.Pool, userID string, filter domain.FilterSearchPagination) ([]domain.OrderUser, error) {
	query := `
		SELECT 
			id, user_id, order_id, product_id, type, price, created_at, updated_at
		FROM 
			order_user
		WHERE 
			user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := db.Query(ctx, query, userID, filter.Pagination.Limit, filter.Pagination.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.OrderUser
	for rows.Next() {
		var order domain.OrderUser
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.OrderID,
			&order.ProductID,
			&order.Type,
			&order.Price,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (q *PaymentQueryImpl) CountOrdersByUserID(ctx context.Context, db *pgxpool.Pool, userID string, filter domain.FilterSearchPagination) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM order_user
		WHERE user_id = $1`

	var count int
	err := db.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}

// Payment operations
func (q *PaymentQueryImpl) CreatePayment(ctx context.Context, tx pgx.Tx, req domain.PaymentHistory) error {
	query := `
		INSERT INTO payment_history (
			id,
			user_id,
			pm_id,
			order_id,
			va_id,
			status,
			expired_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)`

	_, err := tx.Exec(ctx, query,
		req.ID,
		req.UserID,
		req.PaymentId,
		req.OrderID,
		req.VAID,
		req.Status,
		req.ExpiredAt,
	)
	return err
}

func (q *PaymentQueryImpl) GetPaymentByID(ctx context.Context, db *pgxpool.Pool, id string) (domain.PaymentHistory, error) {
	query := `
		SELECT 
			id, user_id, order_id, va_id, status, created_at, updated_at, expired_at
		FROM 
			payment_history
		WHERE 
			order_id = $1`

	var payment domain.PaymentHistory
	err := db.QueryRow(ctx, query, id).Scan(
		&payment.ID,
		&payment.UserID,
		&payment.OrderID,
		&payment.VAID,
		&payment.Status,
		&payment.CreatedAt,
		&payment.UpdatedAt,
		&payment.ExpiredAt,
	)
	return payment, err
}

func (q *PaymentQueryImpl) GetPaymentsByOrderID(ctx context.Context, db *pgxpool.Pool, orderID string) (domain.PaymentHistory, error) {
	query := `
        SELECT 
            id, user_id, order_id, va_id, status, created_at, updated_at, expired_at
        FROM 
            payment_history
        WHERE 
            order_id = $1
        ORDER BY created_at DESC
        LIMIT 1`

	var payment domain.PaymentHistory
	var createdAt, updatedAt, expiredAt time.Time

	err := db.QueryRow(ctx, query, orderID).Scan(
		&payment.ID,
		&payment.UserID,
		&payment.OrderID,
		&payment.VAID,
		&payment.Status,
		&createdAt,
		&updatedAt,
		&expiredAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.PaymentHistory{}, fmt.Errorf("payment with order_id %s not found", orderID)
		}
		return domain.PaymentHistory{}, fmt.Errorf("failed to get payment: %w", err)
	}

	return payment, nil
}

func (q *PaymentQueryImpl) UpdatePaymentStatus(ctx context.Context, tx pgx.Tx, paymentID string, status string) error {
	query := `
		UPDATE payment_history
		SET 
			status = $1,
			updated_at = $2
		WHERE 
			order_id = $3`

	_, err := tx.Exec(ctx, query, status, time.Now(), paymentID)
	return err
}

func (q *PaymentQueryImpl) GetPaymentByOrderID(ctx context.Context, db *pgxpool.Pool, orderID string) (domain.PaymentResponse, error) {
	query := `
        SELECT 
            ph.id, 
            ph.user_id, 
            ph.pm_id as payment_id,
            ph.order_id, 
            ph.va_id,
            ph.status,
            ph.created_at,
            ph.updated_at,
            ph.expired_at,
            p.id as product_id
        FROM 
            payment_history ph
        JOIN 
            order_user ou ON ph.order_id = ou.order_id
        LEFT JOIN 
            product p ON ou.product_id = p.id
        WHERE 
            ph.order_id = $1`

	var payment domain.PaymentResponse
	var createdAt, updatedAt, expiredAt time.Time

	err := db.QueryRow(ctx, query, orderID).Scan(
		&payment.ID,
		&payment.UserID,
		&payment.PaymentId,
		&payment.OrderID,
		&payment.VAID,
		&payment.Status,
		&createdAt,
		&updatedAt,
		&expiredAt,
		&payment.ProductId,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.PaymentResponse{}, fmt.Errorf("payment with order_id %s not found", orderID)
		}
		return domain.PaymentResponse{}, fmt.Errorf("failed to get payment: %w", err)
	}

	// Convert time to Asia/Jakarta using helper
	payment.CreatedAt = helper.ConvertToJakarta(createdAt)
	payment.UpdatedAt = helper.ConvertToJakarta(updatedAt)
	payment.ExpiredAt = helper.ConvertToJakarta(expiredAt)

	return payment, nil
}
