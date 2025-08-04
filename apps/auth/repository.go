package auth

import (
	"context"
	"rent-application/domain"
	"rent-application/internal/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository interface {
	CreateOTP(c context.Context, email string, code string) error
	GetOTP(c context.Context, email string) (string, error)
	CreateUser(c context.Context, user domain.User) error
	GetUserByEmail(c context.Context, email string, level int) (domain.User, error)
}

type authRepository struct {
	db        database.Store
	authQuery AuthQuery
}

func NewAuthRepository(db database.Store, authQuery AuthQuery) AuthRepository {
	return &authRepository{
		db:        db,
		authQuery: authQuery,
	}
}

func (r *authRepository) CreateUser(c context.Context, user domain.User) error {
	return r.db.WithTransaction(c, func(tx pgx.Tx) error {
		err := r.authQuery.CreateUser(c, tx, user)
		if err != nil {
			return err
		}

		err = r.authQuery.CreateWallet(c, tx, user.ID)
		if err != nil {
		}

		return nil
	})
}

func (r *authRepository) GetUserByEmail(c context.Context, email string, level int) (domain.User, error) {
	var user domain.User
	err := r.db.WithoutTransaction(c, func(db *pgxpool.Pool) error {
		var err error
		user, err = r.authQuery.GetUserByEmail(c, db, email, level)
		return err
	})
	return user, err
}

func (r *authRepository) CreateOTP(c context.Context, email string, code string) error {
	return r.db.WithTransaction(c, func(tx pgx.Tx) error {
		return r.authQuery.CreateOTP(c, tx, email, code)
	})
}

func (r *authRepository) GetOTP(c context.Context, email string) (string, error) {
	var otp string
	err := r.db.WithoutTransaction(c, func(db *pgxpool.Pool) error {
		var err error
		otp, err = r.authQuery.GetOTP(c, db, email)
		return err
	})
	return otp, err
}
