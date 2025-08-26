package web_handler

import (
	"github.com/gofiber/fiber/v2"
	"laotop_final/handlers"
	"laotop_final/usecase/customer_service"
	"laotop_final/validation"
)

type CustomerController interface {
	GetCustomerCtr(ctx *fiber.Ctx) error
	CreateCustomerCtr(ctx *fiber.Ctx) error
	UpdateCustomerCtr(ctx *fiber.Ctx) error
	DeleteCustomerCtr(ctx *fiber.Ctx) error
	GetCustomerByIDCtr(ctx *fiber.Ctx) error
	//GetCustomerByCustomerID(c *fiber.Ctx) error
}

type customerController struct {
	cutomerService customer_service.CustomerService
}

func (t *customerController) GetCustomerCtr(ctx *fiber.Ctx) error {
	res, err := t.cutomerService.GetCustomer()
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessResponse(ctx, res)
}
func (t *customerController) CreateCustomerCtr(ctx *fiber.Ctx) error {
	req := new(customer_service.CreateDataCustomerReq)
	if err := ctx.BodyParser(req); err != nil {
		return handlers.NewErrorJson(ctx)
	}

	if errValid := validation.Validate(*req); errValid != nil {
		return handlers.NewErrorValidate(ctx, errValid[0].Error)
	}

	customerID, err := t.cutomerService.CreateCustomer(*req)
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}

	return ctx.JSON(fiber.Map{
		"status":      true,
		"message":     "Create customer success",
		"customer_id": customerID,
	})
}

func (t *customerController) UpdateCustomerCtr(ctx *fiber.Ctx) error {
	req := new(customer_service.UpdateDataCustomerReq)
	err := ctx.BodyParser(req)
	if err != nil {
		return handlers.NewErrorJson(ctx)
	}
	errValid := validation.Validate(*req)
	if errValid != nil {
		return handlers.NewErrorValidate(ctx, errValid[0].Error)
	}
	err = t.cutomerService.UpdateCustomer(*req)
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessMsg(ctx)
}

func (t *customerController) DeleteCustomerCtr(ctx *fiber.Ctx) error {
	req := new(customer_service.DeleteDataCustomerReq)
	err := ctx.BodyParser(req)
	if err != nil {
		return handlers.NewErrorJson(ctx)
	}
	errValid := validation.Validate(*req)
	if errValid != nil {
		return handlers.NewErrorValidate(ctx, errValid[0].Error)
	}
	err = t.cutomerService.DeleteCustomer(*req)
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessMsg(ctx)
}

func (t *customerController) GetCustomerByIDCtr(ctx *fiber.Ctx) error {
	var req customer_service.DataCustomerByIdReq
	if err := ctx.BodyParser(&req); err != nil {
		return handlers.NewErrorResponses(ctx, fiber.NewError(fiber.StatusBadRequest, "Invalid JSON body"))
	}
	res, err := t.cutomerService.GetCustomerByID(req)
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessResponse(ctx, res)
}

func NewCustomerController(
	cutomerService *customer_service.CustomerService,

) CustomerController {
	return &customerController{
		cutomerService: *cutomerService,
	}

}
