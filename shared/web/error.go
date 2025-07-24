package web

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	// Log error jika logger tersedia di context
	if l := c.Locals("logger"); l != nil {
		if logger, ok := l.(*zerolog.Logger); ok {
			logger.Error().Err(err).Int("code", code).Msg("Fiber error handler")
		}
	}

	response := WebResponse{
		Code:    code,
		Status:  false,
		Message: err.Error(),
		Data:    err.Error(),
	}

	return c.Status(code).JSON(response)
}
