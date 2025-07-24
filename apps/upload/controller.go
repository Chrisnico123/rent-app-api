package upload

import (
	"bytes"
	"io"
	"net/http"
	"rent-application/shared/web"

	"github.com/gofiber/fiber/v2"
)

type uploadController struct {
	uploadService UploadService
}

type UploadController interface {
	Route(app *fiber.App)
}

func NewUploadController(service UploadService) UploadController {
	return &uploadController{
		uploadService: service,
	}
}

func (controller *uploadController) Route(apps *fiber.App) {
	app := apps.Group("/upload")
	app.Post("/image", controller.UploadImage)
}

func (controller *uploadController) UploadImage(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return web.ErrBadRequest("file is required")
	}
	f, err := file.Open()
	if err != nil {
		return web.ErrBadRequest("failed to open file")
	}
	defer f.Close()
	// Membaca file ke dalam buffer
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, f); err != nil {
		return web.ErrInternalServer("Failed to read file into buffer")
	}
	url, err := controller.uploadService.UploadImage(c.Context(), buf.Bytes())
	if err != nil {
		return err
	}
	return c.Status(http.StatusOK).JSON(map[string]interface{}{"url": url})
}
