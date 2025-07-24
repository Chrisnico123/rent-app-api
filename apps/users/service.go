package users

import (
	"context"
	"rent-application/domain"
	"rent-application/shared/web"

	"github.com/jackc/pgx/v5"
)

type UsersService interface {
	GetUserById(ctx context.Context, userId string) (web.User, error)
	UpdateUserProfile(ctx context.Context, userId string, req UsersRequest) error
}

type usersService struct {
	usersRepository UsersRepository
}

func NewUsersService(usersRepo UsersRepository) UsersService {
	return &usersService{
		usersRepository: usersRepo,
	}
}

func (u *usersService) UpdateUserProfile(ctx context.Context, userId string, req UsersRequest) error {
	err := req.Validate()
	if err != nil {
		return web.ErrBadRequest(err.Error())
	}

	err = u.usersRepository.UpdateUserProfile(ctx, domain.UserRequest{
		Id:   userId,
		Img:  req.Img,
		Name: req.Name,
	})
	if err != nil {
		return web.ErrInternalServer(err.Error())
	}

	return nil
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

	res := web.User{
		ID:    data.ID,
		Email: data.Email,
		Name:  data.Name,
		Img:   "",
		Level: data.Level,
	}

	if data.Img != nil {
		res.Img = *data.Img
	}

	return res, nil
}
