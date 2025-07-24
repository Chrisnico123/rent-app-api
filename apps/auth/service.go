package auth

import (
	"context"
	"fmt"
	"rent-application/domain"
	"rent-application/internal/email"
	"rent-application/shared/helper"
	"rent-application/shared/web"
)

type AuthService interface {
	SendOtp(ctx context.Context, req OTPReq) error
	VerifyOtp(ctx context.Context, req OTP) (web.AuthResponse, error)
	RegisterUser(ctx context.Context, user User) error
	RegisterSeller(ctx context.Context, user User) error
}

type authService struct {
	authRepository AuthRepository
	emailClient    *email.ResendClient
}

func NewAuthService(authRepo AuthRepository, emailClient *email.ResendClient) AuthService {
	return &authService{
		authRepository: authRepo,
		emailClient:    emailClient,
	}
}

func (s *authService) RegisterSeller(ctx context.Context, user User) error {
	err := user.Validate()
	if err != nil {
		return web.ErrValidateBadRequest(err.Error(), user)
	}

	existingUser, err := s.authRepository.GetUserByEmail(ctx, user.Email, 2)
	if err != nil {
		return err
	}
	if existingUser.ID != "" {
		return web.ErrConflict("user already exists")
	}

	return s.authRepository.CreateUser(ctx, domain.User{
		ID:    helper.GenerateId(),
		Email: user.Email,
		Name:  user.Name,
		Level: 2,
	})
}

func (s *authService) RegisterUser(ctx context.Context, user User) error {
	err := user.Validate()
	if err != nil {
		return web.ErrValidateBadRequest(err.Error(), user)
	}

	existingUser, err := s.authRepository.GetUserByEmail(ctx, user.Email, 1)
	if err != nil {
		return err
	}
	if existingUser.ID != "" {
		return web.ErrConflict("user already exists")
	}

	return s.authRepository.CreateUser(ctx, domain.User{
		ID:    helper.GenerateId(),
		Email: user.Email,
		Name:  user.Name,
		Level: 1,
	})
}

func (s *authService) SendOtp(ctx context.Context, req OTPReq) error {
	err := req.Validate()
	if err != nil {
		return web.ErrValidateBadRequest(err.Error(), req)
	}

	otpCode := email.GenerateOtpCode()

	existingUser, err := s.authRepository.GetUserByEmail(ctx, req.Email, req.Level)
	if err != nil {
		return err
	}

	if existingUser.ID == "" {
		return web.ErrNotFound("user not found")
	}

	err = s.authRepository.CreateOTP(ctx, req.Email, otpCode)
	if err != nil {
		return err
	}

	html := email.OtpEmailTemplate(otpCode)

	subject := "Kode OTP Anda"
	return s.emailClient.SendEmail(req.Email, subject, html)
}

func (s *authService) VerifyOtp(ctx context.Context, req OTP) (web.AuthResponse, error) {
	err := req.Validate()
	if err != nil {
		return web.AuthResponse{}, web.ErrValidateBadRequest(err.Error(), req)
	}

	data, _ := s.authRepository.GetUserByEmail(ctx, req.Email, req.Level)

	if data.ID == "" {
		return web.AuthResponse{}, web.ErrNotFound("user not found")
	}

	otp, err := s.authRepository.GetOTP(ctx, req.Email)
	if err != nil {
		return web.AuthResponse{}, err
	}
	if otp != req.Code {
		return web.AuthResponse{}, web.ErrBadRequest("OTP is invalid")
	}

	token, err := helper.GenerateJwt(data.ID, fmt.Sprintf("%d", data.Level), data.Email)
	if err != nil {
		return web.AuthResponse{}, err
	}

	return web.AuthResponse{
		Id:    data.ID,
		Token: token,
		Email: data.Email,
		Level: data.Level,
	}, nil
}
