package users

import (
	"context"
	"rent-application/domain"
	"rent-application/internal/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersRepository interface {
	GetUserById(ctx context.Context, userId string) (domain.User, error)
	UpdateUserProfile(ctx context.Context, req domain.UserRequest) error
}

type usersRepository struct {
	db         database.Store
	usersQuery UsersQuery
}

func NewUsersRepository(db database.Store, usersQuery UsersQuery) UsersRepository {
	return &usersRepository{
		db:         db,
		usersQuery: usersQuery,
	}
}

func (u *usersRepository) UpdateUserProfile(ctx context.Context, req domain.UserRequest) error {
	return u.db.WithTransaction(ctx, func(tx pgx.Tx) error {
		return u.usersQuery.UpdateUserProfile(ctx, tx, req)
	})
}

// GetUserById implements UsersRepository.
func (u *usersRepository) GetUserById(ctx context.Context, userId string) (domain.User, error) {
	var user domain.User
	err := u.db.WithoutTransaction(ctx, func(db *pgxpool.Pool) error {
		var err error
		user, err = u.usersQuery.GetUserById(ctx, db, userId)
		return err
	})
	return user, err
}
