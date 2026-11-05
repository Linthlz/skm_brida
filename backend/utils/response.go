package utils

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
)

type response struct {
	Success      bool        `json:"success"`
	Status       string      `json:"status"`
	ResponseCode int         `json:"ResponseCode"`
	Message      string      `json:"message,omitempty"`
	Data         interface{} `json:"data,omitempty"`
	Error        string      `json:"error"`
}
type ResponsePaginatedMeta struct {
	Success      bool           `json:"success"`
	Status       string         `json:"status"`
	ResponseCode int            `json:"ResponseCode"`
	Message      string         `json:"message,omitempty"`
	Data         interface{}    `json:"data,omitempty"`
	Error        string         `json:"error"`
	Meta         PaginationMeta `json:"meta"`
}
type PaginationMeta struct {
	Page      int    `json:"page" example:"1"`
	Limit     int    `json:"limit" example:"10"`
	Total     int    `json:"total" example:"100"`
	TotalPage int    `json:"total_page" example:"10"`
	Filter    string `json:"filter" example:"name=Viktor"`
	Sort      string `json:"sort" example:"-id"`
}

func send(c *fiber.Ctx, code int, status, message string, data interface{}, err string) error {
	return c.Status(code).JSON(response{
		Success:      code < 400,
		Status:       status,
		ResponseCode: code,
		Message:      message,
		Data:         data,
		Error:        err,
	})
}

func Success(c *fiber.Ctx, Message string, data interface{}) error {
	return send(c, fiber.StatusOK, "success", Message, data, "")
}
func CreatedWithData(c *fiber.Ctx, Message string, data interface{}) error {
	return send(c, fiber.StatusCreated, "Created", Message, data, "")
}
func SuccessPaginated(c *fiber.Ctx, Message string, data interface{}, meta PaginationMeta) error {
	return c.Status(fiber.StatusOK).JSON(ResponsePaginatedMeta{
		Success:      true,
		Status:       "success",
		ResponseCode: fiber.StatusOK,
		Message:      Message,
		Data:         data,
		Meta:         meta,
	})
}
func NotFoundPaginated(c *fiber.Ctx, Message string, data interface{}, meta PaginationMeta) error {
	return c.Status(fiber.StatusNotFound).JSON(ResponsePaginatedMeta{
		Success:      false,
		Status:       "Not Found",
		ResponseCode: fiber.StatusNotFound,
		Message:      Message,
		Data:         data,
		Meta:         meta,
	})
}
func Created(c *fiber.Ctx, Message string) error {
	return send(c, fiber.StatusCreated, "Created", Message, nil, "")
}
func BadReq(c *fiber.Ctx, Message string, err string) error {
	return send(c, fiber.StatusBadRequest, "Bad Request", Message, nil, err)
}
func NotFound(c *fiber.Ctx, Message string, err string) error {
	return send(c, fiber.StatusNotFound, "Not Found", Message, nil, err)
}
func UnauthorizedAccess(c *fiber.Ctx, Message string, err string) error {
	return send(c, fiber.StatusUnauthorized, "Unauthorized", Message, nil, err)
}
func Forbidden(c *fiber.Ctx, Message string, err string) error {
	return send(c, fiber.StatusForbidden, "Forbidden", Message, nil, err)
}
func Conflict(c *fiber.Ctx, Message string, err string) error {
	return send(c, fiber.StatusConflict, "Conflict", Message, nil, err)
}

// InternalServerError logs the real cause and never forwards it to the client.
func InternalServerError(c *fiber.Ctx, Message string, err string) error {
	if err != "" {
		log.Printf("[error] %s %s: %s: %s", c.Method(), c.Path(), Message, err)
	}
	return send(c, fiber.StatusInternalServerError, "Internal server Error", Message, nil, "internal server error")
}

// AppError is an error whose message is safe to show to API clients.
type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string { return e.Message }

func ErrValidation(msg string) error { return &AppError{Code: fiber.StatusBadRequest, Message: msg} }
func ErrNotFound(msg string) error   { return &AppError{Code: fiber.StatusNotFound, Message: msg} }
func ErrConflict(msg string) error   { return &AppError{Code: fiber.StatusConflict, Message: msg} }
func ErrForbidden(msg string) error  { return &AppError{Code: fiber.StatusForbidden, Message: msg} }

// HandleError maps service errors to HTTP responses; unknown errors become a generic 500.
func HandleError(c *fiber.Ctx, message string, err error) error {
	var appErr *AppError
	if errors.As(err, &appErr) {
		switch appErr.Code {
		case fiber.StatusBadRequest:
			return BadReq(c, message, appErr.Message)
		case fiber.StatusNotFound:
			return NotFound(c, message, appErr.Message)
		case fiber.StatusConflict:
			return Conflict(c, message, appErr.Message)
		case fiber.StatusForbidden:
			return Forbidden(c, message, appErr.Message)
		}
	}
	return InternalServerError(c, message, err.Error())
}
