package payment

import (
	"context"
	"log"
	"rent-application/internal/middleware"
	"rent-application/shared/constants"
	"rent-application/shared/web"
	"strconv"

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

	app.Get("/history",
		middleware.RoleBasedAuth(constants.RolePublicString),
		c.GetPaymentHistory)
	app.Post("/book",
		middleware.RoleBasedAuth(constants.RoleCustomerString),
		c.PaymentBooking)
	app.Get("/book/:order_id",
		middleware.RoleBasedAuth(constants.RoleCustomerString),
		c.GetPaymentBooking)

	app.Get("/credit/history",
		middleware.RoleBasedAuth(constants.RolePublicString),
		c.GetHistoryUser)

	app.Get("/credit",
		middleware.RoleBasedAuth(constants.RolePublicString),
		c.GetBalance)

	app.Post("/credit/topup",
		middleware.RoleBasedAuth(constants.RolePublicString),
		c.TopUpWallet)

	app.Post("/callback/xendit", c.XenditPaymentCallback)
}

func (controller *paymentController) TopUpWallet(c *fiber.Ctx) error {
	var req TopUpRequest
	err := c.BodyParser(&req)
	if err != nil {
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
	data, err := controller.service.CreateTopUP(c.Context(), req, userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "successfully",
		Data:    data,
	})
}

func (controller *paymentController) GetHistoryUser(c *fiber.Ctx) error {
	// Get userID from locals with proper type assertion
	userIDValue := c.Locals("userID")
	if userIDValue == nil {
		return web.ErrBadRequest("userID not found in context")
	}

	userID, ok := userIDValue.(string)
	if !ok || userID == "" {
		return web.ErrBadRequest("invalid userID format")
	}
	data, err := controller.service.GetListTopUP(c.Context(), userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "successfully",
		Data:    data,
	})
}

func (controller *paymentController) GetBalance(c *fiber.Ctx) error {
	// Get userID from locals with proper type assertion
	userIDValue := c.Locals("userID")
	if userIDValue == nil {
		return web.ErrBadRequest("userID not found in context")
	}

	userID, ok := userIDValue.(string)
	if !ok || userID == "" {
		return web.ErrBadRequest("invalid userID format")
	}
	data, err := controller.service.GetBalanceUser(c.Context(), userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "successfully",
		Data:    data,
	})
}

func (controller *paymentController) GetPaymentHistory(c *fiber.Ctx) error {
	var filter web.FIlterPaymentHistory
	if err := c.QueryParser(&filter); err != nil {
		return web.ErrValidateBadRequest(err.Error(), filter)
	}

	// Get userID from context
	userIDValue := c.Locals("userID")
	if userIDValue == nil {
		return web.ErrBadRequest("userID not found in context")
	}

	userID, ok := userIDValue.(string)
	if !ok || userID == "" {
		return web.ErrBadRequest("invalid userID format")
	}

	// Get roleID from context (as int)
	roleValue := c.Locals("roleID")
	if roleValue == nil {
		return web.ErrBadRequest("role user not found")
	}

	roleID, ok := roleValue.(int)
	if !ok {
		return web.ErrBadRequest("invalid role user format")
	}

	// Convert roleID to string for filter
	filter.Level = strconv.Itoa(roleID)
	filter.UserId = userID

	data, count, err := controller.service.GetPaymentListPayment(c.Context(), filter)
	if err != nil {
		return err
	}

	pageInt, _ := strconv.Atoi(filter.Page)

	return c.Status(fiber.StatusOK).JSON(web.WebResponsePagination{
		Code:      fiber.StatusOK,
		Status:    true,
		Page:      pageInt,
		Count:     len(data),
		TotalData: count,
		Message:   "success",
		Data:      data,
	})
}

func (controller *paymentController) GetPaymentBooking(c *fiber.Ctx) error {
	id := c.Params("order_id")

	data, err := controller.service.GetPaymentByOrderId(c.Context(), id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "successfully get payment booking",
		Data:    data,
	})
}

func (controller *paymentController) PaymentBooking(c *fiber.Ctx) error {
	var req BookRequest
	if err := c.BodyParser(&req); err != nil {
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
