package users

import (
	"context"
	"rent-application/shared/web"

	"github.com/jackc/pgx/v5"
)

type UsersService interface {
	GetUserById(ctx context.Context, userId string) (web.User, error)
}

type usersService struct {
	usersRepository UsersRepository
}

func NewUsersService(usersRepo UsersRepository) UsersService {
	return &usersService{
		usersRepository: usersRepo,
	}
}

// GetUserById implements UsersService.
func (u *usersService) GetUserById(ctx context.Context, userId string) (web.User, error) {
	data, err := u.usersRepository.GetUserById(ctx, userId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return web.User{}, web.ErrNotFound("user not found")
		}
		return web.User{}, web.ErrInternalServer(err.Error())
	}

	return web.User(data), nil
}
