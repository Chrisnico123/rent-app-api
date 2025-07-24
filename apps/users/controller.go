package users

import (
	"rent-application/internal/middleware"
	"rent-application/shared/constants"
	"rent-application/shared/web"

	"github.com/gofiber/fiber/v2"
)

type userController struct {
	service UsersService
}

type UserController interface {
	Route(app *fiber.App)
}

func NewUserController(service UsersService) UserController {
	return &userController{
		service: service,
	}
}

// Handler untuk endpoint /health
func (c *userController) Route(apps *fiber.App) {
	app := apps.Group("/user")

	app.Get("/",
		middleware.RoleBasedAuth(constants.RolePublicString),
		c.GetUserData)

	app.Put("/",
		middleware.RoleBasedAuth(constants.RolePublicString),
		c.UpdateUser,
	)
}

func (controller *userController) UpdateUser(c *fiber.Ctx) error {
	var req UsersRequest
	if err := c.BodyParser(&req); err != nil {
		return web.ErrValidateBadRequest(err.Error(), req)
	}

	userIDValue := c.Locals("userID")
	if userIDValue == nil {
		return web.ErrBadRequest("userID not found in context")
	}

	userID, ok := userIDValue.(string)
	if !ok || userID == "" {
		return web.ErrBadRequest("invalid userID format")
	}

	err := controller.service.UpdateUserProfile(c.Context(), userID, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(web.WebResponse{
		Code:    fiber.StatusCreated,
		Status:  true,
		Message: "updated successfully",
	})
}

func (controller *userController) GetUserData(c *fiber.Ctx) error {
	userIDValue := c.Locals("userID")
	if userIDValue == nil {
		return web.ErrBadRequest("userID not found in context")
	}

	userID, ok := userIDValue.(string)
	if !ok || userID == "" {
		return web.ErrBadRequest("invalid userID format")
	}

	data, err := controller.service.GetUserById(c.Context(), userID)
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
