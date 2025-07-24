package payment

import (
	"context"
	"log"
	"rent-application/internal/middleware"
	"rent-application/shared/constants"
	"rent-application/shared/web"

	"github.com/gofiber/fiber/v2"
)

type paymentController struct {
	service PaymentService
}

type PaymentController interface {
	Route(app *fiber.App)
}

func NewPaymentController(service PaymentService) PaymentController {
	return &paymentController{
		service: service,
	}
}

func (c *paymentController) Route(apps *fiber.App) {
	app := apps.Group("/payment")

	// app.Post("/", c.PaymentTopUp)
	app.Post("/book",
		middleware.RoleBasedAuth(constants.RoleCustomerString),
		c.PaymentBooking)
	app.Get("/book/:order_id",
		middleware.RoleBasedAuth(constants.RoleCustomerString),
		c.GetPaymentBooking)

	app.Post("/callback/xendit", c.XenditPaymentCallback)
}

// Handler untuk endpoint /create-payment-mandiri
// func (controller *paymentController) PaymentTopUp(c *fiber.Ctx) error {
// 	va, err := controller.service.CreatePaymentVA(context.Background())

// 	if err != nil {
// 		panic(err)
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"status":  true,
// 		"message": "Payment created successfully",
// 		"data":    va,
// 	})
// }

func (controller *paymentController) GetPaymentBooking(c *fiber.Ctx) error {
	id := c.Params("order_id")

	data, err := controller.service.GetPaymentByOrderId(c.Context(), id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "product updated successfully",
		Data:    data,
	})
}

func (controller *paymentController) PaymentBooking(c *fiber.Ctx) error {
	var req web.BookRequest
	if err := c.QueryParser(&req); err != nil {
		return web.ErrValidateBadRequest(err.Error(), req)
	}

	// Get userID from locals with proper type assertion
	userIDValue := c.Locals("userID")
	if userIDValue == nil {
		return web.ErrBadRequest("userID not found in context")
	}

	userID, ok := userIDValue.(string)
	if !ok || userID == "" {
		return web.ErrBadRequest("invalid userID format")
	}

	resp, err := controller.service.CreateBooking(c.Context(), userID, req)
	if err != nil {
		// Handle the error from service layer
		return err
	}

	// Return the successful response
	return c.JSON(web.WebResponse{
		Code:    fiber.StatusCreated,
		Status:  true,
		Message: "successfully",
		Data:    resp,
	})
}

func (controller *paymentController) XenditPaymentCallback(c *fiber.Ctx) error {
	var callbackPayload web.XenditVACallback

	// Parse payload callback
	if err := c.BodyParser(&callbackPayload); err != nil {
		log.Printf("Error parsing body: %v", err)
		return web.ErrValidateBadRequest("Invalid request body: "+err.Error(), nil)
	}

	// Validate signature header
	callbackSignature := c.Get("x-callback-token")
	if callbackSignature == "" {
		callbackSignature = c.Get("X-Callback-Token")
	}

	if !controller.validateXenditCallback(callbackSignature) {
		return web.ErrUnauthorized("Invalid callback signature")
	}

	err := controller.service.AfterPaymentHandler(context.Background(), callbackPayload)
	if err != nil {
		return web.ErrInternalServer(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
	})
}

// validateXenditCallback memverifikasi signature callback
func (controller *paymentController) validateXenditCallback(signature string) bool {
	// Implementasi validasi signature sesuai dokumentasi Xendit
	// https://developers.xendit.co/api-reference/#callbacks
	expectedSignature := "D9VBAsQkZjffl01g6d6ZvUoLwkfejPbc2aJkNNxxq9MxDwF6"
	return signature == expectedSignature
}
