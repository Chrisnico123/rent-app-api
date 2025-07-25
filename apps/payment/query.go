package payment

import (
	"context"
	"encoding/json"
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
	GetPaymentHistory(ctx context.Context, db *pgxpool.Pool, filter domain.FilterPaymentHistory) ([]domain.PaymentHistoryResponse, error)
	CountPaymentHistory(ctx context.Context, db *pgxpool.Pool, filter domain.FilterPaymentHistory) (int, error)

	// Data Crucial
	CreateEncryptedImg(ctx context.Context, tx pgx.Tx, req domain.EncryptedImg) error
	CreateEncryptedIv(ctx context.Context, tx pgx.Tx, req domain.EncryptedIv) error
	GetDataEncryImg(ctx context.Context, db *pgxpool.Pool, userId string) (domain.DataEncryRes, error)
}

type PaymentQueryImpl struct{}

func NewPaymentQuery() PaymentQuery {
	return &PaymentQueryImpl{}
}

// COuntPaymentHistory implements PaymentQuery.
func (q *PaymentQueryImpl) CountPaymentHistory(ctx context.Context, db *pgxpool.Pool, filter domain.FilterPaymentHistory) (int, error) {
	filterQuery, _ := filter.QueryBuildFilterPayment()

	query := fmt.Sprintf(`
        SELECT 
            COUNT(*)
        FROM 
            payment_history ph
        JOIN 
            order_user ou ON ph.order_id = ou.order_id
        LEFT JOIN 
            product p ON ou.product_id = p.id 
            %s`, filterQuery)

	var count int
	err := db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count payment history: %w", err)
	}

	return count, nil
}

func (q *PaymentQueryImpl) GetPaymentHistory(ctx context.Context, db *pgxpool.Pool, filter domain.FilterPaymentHistory) ([]domain.PaymentHistoryResponse, error) {
	filterQuery, pagination := filter.QueryBuildFilterPayment()

	query := fmt.Sprintf(`
        SELECT 
            ph.order_id,
            p.name AS product_name,
            ou.price,
            p.img,
            ph.status,
            ou.type AS method,
            ph.created_at,
            ph.expired_at
        FROM 
            payment_history ph
        JOIN 
            order_user ou ON ph.order_id = ou.order_id
        LEFT JOIN 
            product p ON ou.product_id = p.id 
            %s
        ORDER BY 
            ph.created_at DESC
        %s`, filterQuery, pagination)

	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query payment history: %w", err)
	}
	defer rows.Close()

	var histories []domain.PaymentHistoryResponse
	for rows.Next() {
		var history domain.PaymentHistoryResponse
		var imgBytes []byte
		var createdAt interface{}
		var expiredAt interface{}
		var price float64

		err := rows.Scan(
			&history.OrderId,
			&history.ProductName,
			&price,
			&imgBytes,
			&history.Status,
			&history.Method,
			&createdAt,
			&expiredAt,
		)
		if err != nil {
			return nil, err
		}

		// Format price as string with 2 decimal places
		history.Price = fmt.Sprintf("%.2f", price)

		// Handle product image (same pattern as ProductResponseList)
		if len(imgBytes) > 0 {
			_ = json.Unmarshal(imgBytes, &history.Img) // Assuming img is []string in PaymentHistoryResponse
		}

		// Handle timestamps (same pattern as ProductResponseList)
		if t, ok := createdAt.(time.Time); ok {
			history.CreatedDate = helper.ConvertToJakarta(t)
		}
		if t, ok := expiredAt.(time.Time); ok {
			history.ExpiredDate = helper.ConvertToJakarta(t)
		}

		histories = append(histories, history)
	}

	if len(histories) == 0 {
		return nil, pgx.ErrNoRows
	}

	return histories, nil
}

// CreateEncryptedImg implements PaymentQuery.
func (q *PaymentQueryImpl) CreateEncryptedImg(ctx context.Context, tx pgx.Tx, req domain.EncryptedImg) error {
	query := `
        INSERT INTO encrypted_images (
            id, 
            user_id, 
            encrypted_url
        ) VALUES (
            $1, $2, $3
        )`

	_, err := tx.Exec(ctx, query,
		req.Id,
		req.UserId,
		req.EncryUrl,
	)
	return err
}

// CreateEncryptedIv implements PaymentQuery.
func (q *PaymentQueryImpl) CreateEncryptedIv(ctx context.Context, tx pgx.Tx, req domain.EncryptedIv) error {
	query := `
        INSERT INTO encryption_ivs (
            id, 
            encry_id, 
            iv
        ) VALUES (
            $1, $2, $3
        )`

	_, err := tx.Exec(ctx, query,
		req.Id,
		req.EncryId,
		req.Ivs,
	)
	return err
}

// GetDataEncryImg implements PaymentQuery.
func (q *PaymentQueryImpl) GetDataEncryImg(ctx context.Context, db *pgxpool.Pool, userId string) (domain.DataEncryRes, error) {
	query := `
        SELECT 
            ei.id,
            ei.user_id,
            ei.encrypted_url,
            eiv.iv
        FROM 
            encrypted_images ei
        JOIN 
            encryption_ivs eiv ON ei.id = eiv.encry_id
        WHERE 
            ei.user_id = $1
        ORDER BY 
            ei.created_at DESC
        LIMIT 1`

	var result domain.DataEncryRes
	err := db.QueryRow(ctx, query, userId).Scan(
		&result.Id,
		&result.UserId,
		&result.EncryUrl,
		&result.Ivs,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.DataEncryRes{}, pgx.ErrNoRows
		}
		return domain.DataEncryRes{}, fmt.Errorf("failed to get encrypted image: %w", err)
	}

	return result, nil
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
			price,
			start_date,
			end_date,
			description
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)`

	_, err := tx.Exec(ctx, query,
		req.ID,
		req.UserID,
		req.OrderID,
		req.ProductID,
		req.Type,
		req.Price,
		req.StartDate,
		req.EndDate,
		req.Desc,
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
			u.username,
			u.email,
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
			users u ON ph.user_id = u.id
        LEFT JOIN 
            product p ON ou.product_id = p.id
        WHERE 
            ph.order_id = $1`

	var payment domain.PaymentResponse
	var createdAt, updatedAt, expiredAt time.Time

	err := db.QueryRow(ctx, query, orderID).Scan(
		&payment.ID,
		&payment.UserID,
		&payment.Username,
		&payment.Email,
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
