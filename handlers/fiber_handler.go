package handlers

import (
	"github.com/gofiber/fiber/v2"
	"laotop_final/errs"
	"net/http"
)

var (
	code    int
	message string
)

type ErrorResponse struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
}

func NewErrorResponses(ctx *fiber.Ctx, err error) error {
	switch e := err.(type) {
	case errs.AppError:
		code = e.Status
		message = e.Message
	case error:
		code = http.StatusUnprocessableEntity
		message = err.Error()
	}
	errorResponse := ErrorResponse{
		Status: false,
		Error:  message,
	}
	return ctx.Status(code).JSON(errorResponse)
}

func NewSuccessResponse(ctx *fiber.Ctx, data interface{}) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"status": true,
		"data":   data,
	})
}
func NewErrorAuthResponse(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"status": false,
		"error":  "Authorization",
	})
}
func NewSuccessMsg(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"data":   TxtSuccess,
	})
}
func NewSuccessResponseSignIn(ctx *fiber.Ctx, data interface{}, token string) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"status":       true,
		"data":         data,
		"access_token": token,
	})
}

func NewErrorPermission(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"status": false,
		"error":  "No Permission",
	})
}
func NewErrorJson(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"status": true,
		"error":  "invalid json format",
	})
}
func NewErrorValidate(ctx *fiber.Ctx, data interface{}) error {
	validateError := fiber.Map{
		"error":  data,
		"status": false,
	}
	return ctx.Status(http.StatusUnprocessableEntity).JSON(validateError)
}

var TxtSuccess = "success"
