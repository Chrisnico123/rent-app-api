package auth

import (
	"rent-application/shared/web"

	"github.com/gofiber/fiber/v2"
)

type authController struct {
	authService AuthService
}

// Depend
type AuthController interface {
	Route(app *fiber.App)
}

// Define
func NewAuthController(authService AuthService) AuthController {
	return &authController{
		authService: authService,
	}
}

func (controller *authController) Route(apps *fiber.App) {
	app := apps.Group("/auth")
	app.Post("/send-otp", controller.SendOtp)
	app.Post("/verify-otp", controller.VerifyOtp)
	app.Post("/register", controller.RegisterUser)
	app.Post("/register-seller", controller.RegisterSeller)
}

func (controller *authController) RegisterSeller(c *fiber.Ctx) error {
	var user User
	err := c.BodyParser(&user)
	if err != nil {
		return web.ErrValidateBadRequest(err.Error(), user)
	}

	err = controller.authService.RegisterSeller(c.Context(), user)
	if err != nil {
		return err
	}

	return c.JSON(web.WebResponse{
		Code:    fiber.StatusCreated,
		Status:  true,
		Message: "seller registered successfully",
	})
}

// Handler untuk mengirim OTP ke email
func (controller *authController) SendOtp(c *fiber.Ctx) error {
	var req OTPReq

	err := c.BodyParser(&req)
	if err != nil {
		return web.ErrValidateBadRequest(err.Error(), req)
	}

	err = controller.authService.SendOtp(c.Context(), req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "success",
	})
}

// Handler untuk verifikasi OTP
func (controller *authController) VerifyOtp(c *fiber.Ctx) error {
	var req OTP
	err := c.BodyParser(&req)
	if err != nil {
		return web.ErrValidateBadRequest(err.Error(), req)
	}
	data, err := controller.authService.VerifyOtp(c.Context(), req)
	if err != nil {
		return err
	}

	return c.JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "success",
		Data:    data,
	})
}

func (controller *authController) RegisterUser(c *fiber.Ctx) error {
	var user User
	err := c.BodyParser(&user)
	if err != nil {
		return web.ErrValidateBadRequest(err.Error(), user)
	}

	err = controller.authService.RegisterUser(c.Context(), user)
	if err != nil {
		return err
	}

	return c.JSON(web.WebResponse{
		Code:    fiber.StatusCreated,
		Status:  true,
		Message: "user registered successfully",
	})
}
