package users

import (
	"context"
	"rent-application/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersQuery interface {
	GetUserById(c context.Context, db *pgxpool.Pool, userId string) (domain.User, error)
	UpdateUserProfile(c context.Context, tx pgx.Tx, req domain.UserRequest) error
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

	query = `SELECT id, email, username, level, img FROM users WHERE id = $1`
	err = db.QueryRow(c, query, userId).Scan(&user.ID, &user.Email, &user.Name, &user.Level, &user.Img)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.User{}, pgx.ErrNoRows
		}
		return domain.User{}, err
	}
	return user, nil
}

// UpdateUserImgAndUsername implements UsersQuery.
func (u *UsersQueryImpl) UpdateUserProfile(c context.Context, tx pgx.Tx, req domain.UserRequest) error {
	query := `
		UPDATE users 
		SET img = $1, 
		    username = $2
		WHERE id = $3
	`

	_, err := tx.Exec(c, query, req.Img, req.Name, req.Id)
	if err != nil {
		return err
	}
	return nil
}
