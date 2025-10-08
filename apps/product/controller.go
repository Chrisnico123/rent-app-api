package product

import (
	"rent-application/internal/middleware"
	"rent-application/shared/constants"
	"rent-application/shared/web"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type categoryController struct {
	productService ProductService
}

type CategoryController interface {
	Route(app *fiber.App)
}

func NewCategoryController(productService ProductService) CategoryController {
	return &categoryController{
		productService: productService,
	}
}

func (controller *categoryController) Route(apps *fiber.App) {
	app := apps.Group("/product")

	// Category routes
	app.Post("/category",
		middleware.RoleBasedAuth(constants.RoleAdminString),
		controller.CreateCategory)
	app.Put("/category/:id",
		middleware.RoleBasedAuth(constants.RoleAdminString),
		controller.UpdateCategory)
	app.Delete("/category/:id",
		middleware.RoleBasedAuth(constants.RoleAdminString),
		controller.DeleteCategory)

	app.Get("/category", controller.ListCategory)
	app.Get("/category/:id", controller.GetCategoryByID)

	// Product routes
	app.Post("/",
		middleware.RoleBasedAuth(constants.RoleSellerString),
		controller.CreateProduct)
	app.Put("/:id",
		middleware.RoleBasedAuth(constants.RoleSellerString),
		controller.UpdateProduct)
	app.Delete("/:id",
		middleware.RoleBasedAuth(constants.RoleSellerString),
		controller.DeleteProduct)
	app.Get("/book/:id",
		middleware.RoleBasedAuth(constants.RoleCustomerString),
		controller.GetDataBookProduct)

	app.Get("/:id", controller.GetProductByID)
	app.Get("/", controller.ListProduct)
	app.Get("/list/seller",
		middleware.RoleBasedAuth(constants.RoleSellerString),
		controller.ListProductSeller)

	// Banner routes
	app.Post("/banner",
		middleware.RoleBasedAuth(constants.RoleAdminString),
		controller.CreateBanner)
	app.Put("/banner/:id",
		middleware.RoleBasedAuth(constants.RoleAdminString),
		controller.UpdateBanner)
	app.Delete("/banner/:id",
		middleware.RoleBasedAuth(constants.RoleAdminString),
		controller.DeleteBanner)

	app.Get("/banner/list", controller.ListBanner)
	app.Get("/banner/:id", controller.GetBannerById)
}

func (controller *categoryController) GetDataBookProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.ErrBadRequest("banner ID cannot be empty")
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

	resp, err := controller.productService.GetProductBookById(c.Context(), id, userID)
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

func (controller *categoryController) CreateBanner(c *fiber.Ctx) error {
	var req BannerRequest
	if err := c.BodyParser(&req); err != nil {
		return web.ErrValidateBadRequest(err.Error(), req)
	}

	err := controller.productService.CreateBanner(c.Context(), req)
	if err != nil {
		return web.ErrInternalServer(err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(web.WebResponse{
		Code:    fiber.StatusCreated,
		Status:  true,
		Message: "banner created successfully",
	})
}

func (controller *categoryController) UpdateBanner(c *fiber.Ctx) error {
	var req BannerRequest
	id := c.Params("id")
	if id == "" {
		return web.ErrBadRequest("banner ID cannot be empty")
	}

	if err := c.BodyParser(&req); err != nil {
		return web.ErrValidateBadRequest(err.Error(), req)
	}

	err := controller.productService.UpdateBanner(c.Context(), id, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "banner updated successfully",
	})
}

func (controller *categoryController) DeleteBanner(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.ErrBadRequest("banner ID cannot be empty")
	}

	err := controller.productService.DeleteBanner(c.Context(), id)
	if err != nil {
		return web.ErrInternalServer(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "banner deleted successfully",
	})
}

func (controller *categoryController) GetBannerById(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.ErrBadRequest("banner ID cannot be empty")
	}

	data, err := controller.productService.GetBannerById(c.Context(), id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "success",
		Data:    data,
	})
}

func (controller *categoryController) ListBanner(c *fiber.Ctx) error {
	var filter web.FilterBannerPagination
	if err := c.QueryParser(&filter); err != nil {
		return web.ErrValidateBadRequest(err.Error(), filter)
	}

	result, count, err := controller.productService.GetListBanner(c.Context(), filter)
	if err != nil {
		return err
	}

	pageInt, _ := strconv.Atoi(filter.Page)

	return c.Status(fiber.StatusOK).JSON(web.WebResponsePagination{
		Code:      fiber.StatusOK,
		Status:    true,
		Page:      pageInt,
		Count:     len(result),
		TotalData: count,
		Message:   "success",
		Data:      result,
	})
}

func (controller *categoryController) ListProductSeller(c *fiber.Ctx) error {
	var filter web.FilterProduct
	if err := c.QueryParser(&filter); err != nil {
		return web.ErrValidateBadRequest(err.Error(), filter)
	}

	roleValue := c.Locals("roleID")

	roleID, _ := roleValue.(int)

	if roleID == constants.SellerRole {
		userIDValue := c.Locals("userID")
		if userIDValue == nil {
			return web.ErrBadRequest("userID not found in context")
		}
		userID, ok := userIDValue.(string)
		if !ok || userID == "" {
			return web.ErrBadRequest("invalid userID format")
		}
		filter.UserId = userID
	}

	result, count, err := controller.productService.GetAllProducts(c.Context(), filter)
	if err != nil {
		return web.ErrInternalServer(err.Error())
	}

	pageInt, _ := strconv.Atoi(filter.Page)

	return c.Status(fiber.StatusOK).JSON(web.WebResponsePagination{
		Code:      fiber.StatusOK,
		Status:    true,
		Page:      pageInt,
		Count:     len(result),
		TotalData: count,
		Message:   "success",
		Data:      result,
	})
}

func (controller *categoryController) ListProduct(c *fiber.Ctx) error {
	var filter web.FilterProduct
	if err := c.QueryParser(&filter); err != nil {
		return web.ErrValidateBadRequest(err.Error(), filter)
	}
	result, count, err := controller.productService.GetAllProducts(c.Context(), filter)
	if err != nil {
		return web.ErrInternalServer(err.Error())
	}

	pageInt, _ := strconv.Atoi(filter.Page)

	return c.Status(fiber.StatusOK).JSON(web.WebResponsePagination{
		Code:      fiber.StatusOK,
		Status:    true,
		Page:      pageInt,
		Count:     len(result),
		TotalData: count,
		Message:   "success",
		Data:      result,
	})
}

func (controller *categoryController) CreateProduct(c *fiber.Ctx) error {
	var req Product
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
	req.SellerID = userID
	err := controller.productService.CreateProduct(c.Context(), req)
	if err != nil {
		return web.ErrInternalServer(err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(web.WebResponse{
		Code:    fiber.StatusCreated,
		Status:  true,
		Message: "product created successfully",
	})
}

func (controller *categoryController) UpdateProduct(c *fiber.Ctx) error {
	var req Product
	id := c.Params("id")
	if err := c.BodyParser(&req); err != nil {
		return web.ErrValidateBadRequest(err.Error(), req)
	}

	req.ID = id
	err := controller.productService.UpdateProduct(c.Context(), req)
	if err != nil {
		return web.ErrInternalServer(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "product updated successfully",
	})
}

func (controller *categoryController) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.ErrBadRequest("product ID cannot be empty")
	}

	err := controller.productService.SoftDeleteProduct(c.Context(), id)
	if err != nil {
		return web.ErrInternalServer(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "product deleted successfully",
	})
}

func (controller *categoryController) GetProductByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.ErrBadRequest("product ID cannot be empty")
	}

	data, err := controller.productService.GetProductByID(c.Context(), id)
	if err != nil {
		return web.ErrInternalServer(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "success",
		Data:    data,
	})
}

func (controller *categoryController) GetCategoryByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.ErrBadRequest("category ID cannot be empty")
	}

	data, err := controller.productService.GetCategoryByID(c.Context(), id)
	if err != nil {
		return web.ErrInternalServer(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "success",
		Data:    data,
	})
}

func (controller *categoryController) CreateCategory(c *fiber.Ctx) error {
	var req Category
	if err := c.BodyParser(&req); err != nil {
		return web.ErrValidateBadRequest(err.Error(), req)
	}
	err := controller.productService.CreateCategory(c.Context(), req)
	if err != nil {
		return web.ErrInternalServer(err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(web.WebResponse{
		Code:    fiber.StatusCreated,
		Status:  true,
		Message: "category created successfully",
	})
}

func (controller *categoryController) UpdateCategory(c *fiber.Ctx) error {
	var req Category
	id := c.Params("id")
	if err := c.BodyParser(&req); err != nil {
		return web.ErrValidateBadRequest(err.Error(), req)
	}

	req.ID = id
	err := controller.productService.UpdateCategory(c.Context(), req)
	if err != nil {
		return web.ErrInternalServer(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "category updated successfully",
	})
}

func (controller *categoryController) DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	err := controller.productService.SoftDeleteCategory(c.Context(), id)
	if err != nil {
		return web.ErrInternalServer(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "category deleted successfully",
	})
}

func (controller *categoryController) ListCategory(c *fiber.Ctx) error {
	var filter web.FilterSearchPagination
	if err := c.QueryParser(&filter); err != nil {
		return web.ErrValidateBadRequest(err.Error(), filter)
	}
	result, count, err := controller.productService.GetAllCategories(c.Context(), filter)
	if err != nil {
		return web.ErrInternalServer(err.Error())
	}

	pageInt, _ := strconv.Atoi(filter.Page)

	return c.Status(fiber.StatusOK).JSON(web.WebResponsePagination{
		Code:      fiber.StatusOK,
		Status:    true,
		Page:      pageInt,
		Count:     len(result),
		TotalData: count,
		Message:   "success",
		Data:      result,
	})
}
