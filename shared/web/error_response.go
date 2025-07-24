package web

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ErrConflict(message string) *fiber.Error {
	return fiber.NewError(fiber.StatusConflict, message)
}

func ErrBadRequest(message string) *fiber.Error {
	return fiber.NewError(fiber.StatusBadRequest, message)
}

func ErrUnauthorized(message string) *fiber.Error {
	return fiber.NewError(fiber.StatusUnauthorized, message)
}

func ErrForbidden(message string) *fiber.Error {
	return fiber.NewError(fiber.StatusForbidden, message)
}

func ErrNotFound(message string) *fiber.Error {
	return fiber.NewError(fiber.StatusNotFound, message)
}

func ErrUnprocessableEntity(message string) *fiber.Error {
	return fiber.NewError(fiber.StatusUnprocessableEntity, message)
}

func ErrInternalServer(message string) *fiber.Error {
	return fiber.NewError(fiber.StatusInternalServerError, message)
}

func ErrNotImplemented(message string) *fiber.Error {
	return fiber.NewError(fiber.StatusNotImplemented, message)
}

func ErrServiceUnavailable(message string) *fiber.Error {
	return fiber.NewError(fiber.StatusServiceUnavailable, message)
}

func ErrTooManyRequests(message string) *fiber.Error {
	return fiber.NewError(fiber.StatusTooManyRequests, message)
}

// ErrValidateBadRequest is a function to handle error validation
func ErrValidateBadRequest(message string, data interface{}) *fiber.Error {
	resMessage := message
	if strings.Contains(message, "cannot unmarshal") && strings.Contains(message, " of type ") {
		split := strings.Split(message, ".")
		values := strings.Split(split[len(split)-1], " of type ")
		resMessage = fmt.Sprintf("Field '%s' must be filled with a value of type %s", values[0], values[1])
	} else if strings.Contains(message, "invalid character") {
		resMessage = "Bad body request, check the JSON formatting"
	} else if strings.Contains(message, "Error:Field validation") && strings.Contains(message, "' failed on the '") {
		values := strings.Split(strings.Split(message, "validation for '")[1], "' failed on the '")
		t := reflect.TypeOf(data)
		field, _ := t.FieldByName(values[0])
		validate := field.Tag.Get("validate")
		tag := field.Tag.Get("json")

		// check the validation error further
		switch {
		case strings.Contains(message, "'required'"):
			resMessage = fmt.Sprintf("Field '%s' must be filled", tag)
		case strings.Contains(message, "'max'"):
			maxChar := strings.Split(strings.Split(validate, "max=")[1], ",")[0]
			prefix := "character"
			suffix := " character"
			if field.Type.Kind() == reflect.Int || field.Type.Kind() == reflect.Ptr {
				prefix = "number"
				suffix = " int"
			}
			resMessage = fmt.Sprintf("Field '%s' exceeded the maximum %s limit of: %s %s", tag, prefix, maxChar, suffix)
		case strings.Contains(message, "'min'"):
			prefix := "character"
			suffix := "character"
			if field.Type.Kind() == reflect.Int || field.Type.Kind() == reflect.Ptr {
				prefix = "number"
				suffix = "int"
			}
			minChar := strings.Split(strings.Split(validate, "min=")[1], ",")[0]
			resMessage = fmt.Sprintf("Field '%s' exceeded the minimum %s limit of: %s %s", tag, prefix, minChar, suffix)
		case strings.Contains(message, "failed on the"):
			format := strings.Split(strings.Split(message, "failed on the '")[1], "' tag")[0]
			resMessage = fmt.Sprintf("Field '%s' must have a '%s' format", tag, format)
		default:
			resMessage = "There's error on body requests"
		}
	}
	return ErrBadRequest(resMessage)
}
