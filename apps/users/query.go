package users

import (
	"context"
	"rent-application/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersQuery interface {
	GetUserById(c context.Context, db *pgxpool.Pool, userId string) (domain.User, error)
}

type UsersQueryImpl struct{}

func NewUsersQuery() UsersQuery {
	return &UsersQueryImpl{}
}

// GetUserById implements UsersQuery.
func (u *UsersQueryImpl) GetUserById(c context.Context, db *pgxpool.Pool, userId string) (domain.User, error) {
	var user domain.User
	var query string
	var err error

	query = `SELECT id, email, username, level FROM users WHERE id = $1`
	err = db.QueryRow(c, query, userId).Scan(&user.ID, &user.Email, &user.Name, &user.Level)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.User{}, pgx.ErrNoRows
		}
		return domain.User{}, err
	}
	return user, nil
}
