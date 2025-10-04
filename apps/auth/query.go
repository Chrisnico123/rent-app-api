package auth

import (
	"context"
	"rent-application/domain"
	"rent-application/shared/helper"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthQuery interface {
	CreateUser(c context.Context, tx pgx.Tx, req domain.User) error
	CreateOTP(c context.Context, tx pgx.Tx, email string, code string) error
	DeleteOTP(c context.Context, tx pgx.Tx, email string) error
	CreateWallet(ctx context.Context, tx pgx.Tx, userId string) error

	GetUserByEmail(c context.Context, db *pgxpool.Pool, email string, level int) (domain.User, error)
	GetOTP(c context.Context, db *pgxpool.Pool, email string) (string, error)
}

type AuthQueryImpl struct{}

func NewAuthQuery() AuthQuery {
	return &AuthQueryImpl{}
}

// CreateWallet implements PaymentQuery.
func (q *AuthQueryImpl) CreateWallet(ctx context.Context, tx pgx.Tx, userId string) error {
	query := `
		INSERT INTO wallets (
			id, 
			user_id
		) VALUES (
			$1, $2
		)`

	_, err := tx.Exec(ctx, query,
		helper.GenerateId(),
		userId,
	)

	return err
}

func (q *AuthQueryImpl) CreateUser(c context.Context, tx pgx.Tx, req domain.User) error {
	query := `INSERT INTO users (id, email, username, level) VALUES ($1, $2, $3, $4)`
	_, err := tx.Exec(c, query, req.ID, req.Email, req.Name, req.Level)
	return err
}

func (q *AuthQueryImpl) GetUserByEmail(c context.Context, db *pgxpool.Pool, email string, level int) (domain.User, error) {
	var user domain.User
	var query string
	var err error

	query = `SELECT id, email, username, level FROM users WHERE email = $1 AND level = $2`
	err = db.QueryRow(c, query, email, level).Scan(&user.ID, &user.Email, &user.Name, &user.Level)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.User{}, nil
		}
		return domain.User{}, err
	}
	return user, nil
}

func (q *AuthQueryImpl) CreateOTP(c context.Context, tx pgx.Tx, email string, code string) error {
	query := `INSERT INTO otp_codes (email, otp_code) VALUES ($1, $2)`
	_, err := tx.Exec(c, query, email, code)
	return err
}

func (q *AuthQueryImpl) GetOTP(c context.Context, db *pgxpool.Pool, email string) (string, error) {
	var otp string
	query := `SELECT otp_code FROM otp_codes WHERE email = $1 AND remove_at > NOW() ORDER BY created_at DESC LIMIT 1`
	err := db.QueryRow(c, query, email).Scan(&otp)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return otp, nil
}

func (q *AuthQueryImpl) DeleteOTP(c context.Context, tx pgx.Tx, email string) error {
	query := `DELETE FROM otp_codes WHERE email = $1`
	_, err := tx.Exec(c, query, email)
	return err
}
