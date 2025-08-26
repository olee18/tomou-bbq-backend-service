package web_handler

import (
	"github.com/gofiber/fiber/v2"
	"laotop_final/handlers"
	"laotop_final/usecase/bill_service"
	"laotop_final/validation"
)

type BillsController interface {
	GetBillsCtr(ctx *fiber.Ctx) error
	GetBillHistroy(ctx *fiber.Ctx) error
	InsertBilssCtr(ctx *fiber.Ctx) error
	DeleteBillCtr(ctx *fiber.Ctx) error
	UpdateBillCtr(ctx *fiber.Ctx) error
	GetBillitem(ctx *fiber.Ctx) error
	DeleteBillItemCtr(ctx *fiber.Ctx) error
	ShowBillClientByCustomerIDCtr(ctx *fiber.Ctx) error
}
type billsController struct {
	billService bill_service.BillsService
}

func (b *billsController) UpdateBillCtr(ctx *fiber.Ctx) error {
	req := new(bill_service.UpdateBillReq)
	err := ctx.BodyParser(req)
	if err != nil {
		return handlers.NewErrorJson(ctx)
	}
	errValid := validation.Validate(*req)
	if errValid != nil {
		return handlers.NewErrorValidate(ctx, errValid[0].Error)
	}
	err = b.billService.UpdateBIllService(*req)
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessMsg(ctx)
}

func (b *billsController) GetBillHistroy(ctx *fiber.Ctx) error {
	var req bill_service.BillsReq
	if err := ctx.BodyParser(&req); err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	res, err := b.billService.GetBillHistory(req)
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessResponse(ctx, res)
}

func (b *billsController) GetBillitem(ctx *fiber.Ctx) error {
	res, err := b.billService.GetAllBillItemsService()
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessResponse(ctx, res)
}

func (b *billsController) GetBillsCtr(ctx *fiber.Ctx) error {
	var req bill_service.BillsReq
	if err := ctx.BodyParser(&req); err != nil {
		return handlers.NewErrorResponses(ctx, fiber.NewError(fiber.StatusBadRequest, "Invalid JSON body"))
	}
	res, err := b.billService.GetBillHistoryByCustomerID(req)
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessResponse(ctx, res)
}

func (b *billsController) InsertBilssCtr(ctx *fiber.Ctx) error {
	req := new(bill_service.BillsReq)
	err := ctx.BodyParser(req)
	if err != nil {
		return handlers.NewErrorJson(ctx)
	}
	errValid := validation.Validate(*req)
	if errValid != nil {
		return handlers.NewErrorValidate(ctx, errValid[0].Error)
	}
	res, err := b.billService.GenerateBillForCustomerID(*req)
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessResponse(ctx, res)
}
func (b *billsController) ShowBillClientByCustomerIDCtr(ctx *fiber.Ctx) error {
	var req struct {
		CustomerID uint `json:"customer_id"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return handlers.NewErrorResponses(ctx, fiber.NewError(fiber.StatusBadRequest, "invalid request body"))
	}

	if req.CustomerID == 0 {
		return handlers.NewErrorResponses(ctx, fiber.NewError(fiber.StatusBadRequest, "customer_id is required"))
	}

	res, err := b.billService.ShowBillClientByCustomerID(req.CustomerID)
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}

	return handlers.NewSuccessResponse(ctx, res)
}
func (b *billsController) DeleteBillCtr(ctx *fiber.Ctx) error {
	req := new(bill_service.BillDeleteReq)
	err := ctx.BodyParser(req)
	if err != nil {
		return handlers.NewErrorJson(ctx)
	}
	errValid := validation.Validate(*req)
	if errValid != nil {
		return handlers.NewErrorValidate(ctx, errValid[0].Error)
	}
	err = b.billService.DeleteBillService(*req)
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessMsg(ctx)
}
func (b *billsController) DeleteBillItemCtr(ctx *fiber.Ctx) error {
	req := new(bill_service.BillitemDeleteReq)
	err := ctx.BodyParser(req)
	if err != nil {
		return handlers.NewErrorJson(ctx)
	}
	errValid := validation.Validate(*req)
	if errValid != nil {
		return handlers.NewErrorValidate(ctx, errValid[0].Error)
	}
	err = b.billService.DeleteBillitemService(*req)
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessMsg(ctx)
}

func NewBillsController(
	billService *bill_service.BillsService,

) BillsController {
	return &billsController{
		billService: *billService,
	}

}
