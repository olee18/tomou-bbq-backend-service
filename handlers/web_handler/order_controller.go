package web_handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"laotop_final/handlers"
	"laotop_final/usecase/order_service"
	"laotop_final/validation"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type OrderController interface {
	GetOrderCtr(ctx *fiber.Ctx) error
	GetORderByCustomerID(ctx *fiber.Ctx) error
	CreateOrderCtr(ctx *fiber.Ctx) error
	UpdateOrderCtr(ctx *fiber.Ctx) error
	DeleteOrderCtr(ctx *fiber.Ctx) error
	GetCategoryCtr(ctx *fiber.Ctx) error
	GetCategoryByIDCtr(ctx *fiber.Ctx) error
	CreateCategoryCtr(ctx *fiber.Ctx) error
	UpdateCategoryCtr(ctx *fiber.Ctx) error
	DeleteCategoryCtr(ctx *fiber.Ctx) error
	GetMenuByCategoryIDCtr(ctx *fiber.Ctx) error
	GetMenuByIDCtr(ctx *fiber.Ctx) error
	GetMenuCtr(ctx *fiber.Ctx) error
	CreateMenuCtr(ctx *fiber.Ctx) error
	UpdateMenuCtr(ctx *fiber.Ctx) error
	DeleteMenuCtr(ctx *fiber.Ctx) error
	GetAllOrderItemsCtr(ctx *fiber.Ctx) error
	GetAllOrdersNoFilterCtr(ctx *fiber.Ctx) error
	DeleteOrderItems(ctx *fiber.Ctx) error
}

type orderController struct {
	orderService order_service.OrderService
}

func (o *orderController) GetAllOrderItemsCtr(ctx *fiber.Ctx) error {
	items, err := o.orderService.GetAllOrderItems()
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessResponse(ctx, items)
}
func (o *orderController) UpdateOrderStatusCtr(ctx *fiber.Ctx) error {
	req := new(order_service.OrderItemUpdateStatusReq)
	if err := ctx.BodyParser(req); err != nil {
		return handlers.NewErrorJson(ctx)
	}
	if errValid := validation.Validate(*req); errValid != nil {
		return handlers.NewErrorValidate(ctx, errValid[0].Error)
	}
	err := o.orderService.UpdateOrderStatus(*req)
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessMsg(ctx)
}

func (o *orderController) GetORderByCustomerID(ctx *fiber.Ctx) error {
	var req order_service.OrderReq
	if err := ctx.BodyParser(&req); err != nil {
		return handlers.NewErrorResponses(ctx, fiber.NewError(fiber.StatusBadRequest, "Invalid JSON body"))
	}
	res, err := o.orderService.GetOrderByCustomerID(req)
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessResponse(ctx, res)
}

func (o *orderController) CreateOrderCtr(ctx *fiber.Ctx) error {
	req := new(order_service.CreateOrder)
	if err := ctx.BodyParser(req); err != nil {
		return handlers.NewErrorJson(ctx)
	}
	errValid := validation.Validate(*req)
	if errValid != nil {
		return handlers.NewErrorValidate(ctx, errValid[0].Error)
	}

	orderID, err := o.orderService.CreateOrder(*req)
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}

	return ctx.JSON(fiber.Map{
		"status":   true,
		"message":  "Create order success",
		"order_id": orderID,
	})
}

func (o *orderController) GetMenuByCategoryIDCtr(ctx *fiber.Ctx) error {
	var requestData struct {
		CategoryID int `json:"category_id"`
	}
	if err := ctx.BodyParser(&requestData); err != nil {
		return handlers.NewErrorJson(ctx)
	}
	id := requestData.CategoryID
	res, err := o.orderService.GetMenuByCategoryID(uint(id))
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessResponse(ctx, res)
}

func (o *orderController) GetOrderCtr(ctx *fiber.Ctx) error {
	res, err := o.orderService.GetOrderHistory()
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessResponse(ctx, res)
}

func (o *orderController) UpdateOrderCtr(ctx *fiber.Ctx) error {
	req := new(order_service.UpdateOrderReq)
	err := ctx.BodyParser(req)
	if err != nil {
		return handlers.NewErrorJson(ctx)
	}
	errValid := validation.Validate(*req)
	if errValid != nil {
		return handlers.NewErrorValidate(ctx, errValid[0].Error)
	}
	err = o.orderService.UpdateOrder(*req)
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessMsg(ctx)
}

func (o *orderController) DeleteOrderCtr(ctx *fiber.Ctx) error {
	req := new(order_service.DeleteOrderReq)
	err := ctx.BodyParser(req)
	if err != nil {
		return handlers.NewErrorJson(ctx)
	}
	errValid := validation.Validate(*req)
	if errValid != nil {
		return handlers.NewErrorValidate(ctx, errValid[0].Error)
	}
	err = o.orderService.DeleteOrder(*req)
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessMsg(ctx)
}

func (o *orderController) GetCategoryByIDCtr(ctx *fiber.Ctx) error {
	var requestData struct {
		ID int `json:"id"`
	}
	if err := ctx.BodyParser(&requestData); err != nil {
		return handlers.NewErrorJson(ctx)
	}
	id := requestData.ID
	res, err := o.orderService.GetCategoryByID(id)
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessResponse(ctx, res)

}
func (o *orderController) GetMenuByIDCtr(ctx *fiber.Ctx) error {
	var requestData struct {
		ID int `json:"id"`
	}
	if err := ctx.BodyParser(&requestData); err != nil {
		return handlers.NewErrorJson(ctx)
	}
	id := requestData.ID
	res, err := o.orderService.GetMenuByID(id)
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessResponse(ctx, res)
}

func (o *orderController) GetCategoryCtr(ctx *fiber.Ctx) error {
	res, err := o.orderService.GetCategory()
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessResponse(ctx, res)
}

func (o *orderController) CreateCategoryCtr(ctx *fiber.Ctx) error {
	req := new(order_service.CreateCategoryReq)
	err := ctx.BodyParser(req)
	if err != nil {
		return handlers.NewErrorJson(ctx)
	}
	errValid := validation.Validate(*req)
	if errValid != nil {
		return handlers.NewErrorValidate(ctx, errValid[0].Error)
	}
	err = o.orderService.CreateCategory(*req)
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessMsg(ctx)
}

func (o *orderController) UpdateCategoryCtr(ctx *fiber.Ctx) error {
	req := new(order_service.UpdateCategoryReq)
	err := ctx.BodyParser(req)
	if err != nil {
		return handlers.NewErrorJson(ctx)
	}
	errValid := validation.Validate(*req)
	if errValid != nil {
		return handlers.NewErrorValidate(ctx, errValid[0].Error)
	}
	err = o.orderService.UpdateCategory(*req)
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessMsg(ctx)
}

func (o *orderController) DeleteCategoryCtr(ctx *fiber.Ctx) error {
	req := new(order_service.DeleteCategoryReq)
	err := ctx.BodyParser(req)
	if err != nil {
		return handlers.NewErrorJson(ctx)
	}
	errValid := validation.Validate(*req)
	if errValid != nil {
		return handlers.NewErrorValidate(ctx, errValid[0].Error)
	}
	err = o.orderService.DeleteCategory(*req)
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessMsg(ctx)
}

func (o *orderController) GetMenuCtr(ctx *fiber.Ctx) error {
	res, err := o.orderService.GetMenu()
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessResponse(ctx, res)
}

func (o *orderController) CreateMenuCtr(ctx *fiber.Ctx) error {
	// multi-parts file handling
	name := ctx.FormValue("name")
	categoryIDStr := ctx.FormValue("category_id")
	priceStr := ctx.FormValue("price")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		return handlers.NewErrorJson(ctx)
	}
	price, err := strconv.Atoi(priceStr)
	if err != nil {
		return handlers.NewErrorJson(ctx)
	}
	file, err := ctx.FormFile("image")
	if err != nil {
		return handlers.NewErrorJson(ctx)
	}
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
	savePath := fmt.Sprintf("./Asset/image/%s", filename)
	// บันทึกไฟล์ลง local
	if err = ctx.SaveFile(file, savePath); err != nil {
		return handlers.NewErrorJson(ctx)
	}
	//note for allow client read file in this project
	req := order_service.CreateMenuReq{
		CategoryID: uint(categoryID),
		Name:       name,
		Price:      price,
		Image:      filename,
		Status:     true,
	}
	err = o.orderService.CreateMenu(req)
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessMsg(ctx)
}

func (o *orderController) UpdateMenuCtr(ctx *fiber.Ctx) error {
	idStr := ctx.FormValue("id")
	name := ctx.FormValue("name")
	categoryIDStr := ctx.FormValue("category_id")
	priceStr := ctx.FormValue("price")
	statusStr := ctx.FormValue("status")
	currentImage := ctx.FormValue("current_image")

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return handlers.NewErrorJson(ctx)
	}
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil || categoryID <= 0 {
		return handlers.NewErrorJson(ctx)
	}
	price, err := strconv.Atoi(priceStr)
	if err != nil || price < 0 {
		return handlers.NewErrorJson(ctx)
	}
	status, err := strconv.ParseBool(statusStr)
	if err != nil {
		return handlers.NewErrorJson(ctx)
	}

	filename := currentImage

	if file, err := ctx.FormFile("image"); err == nil {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
			return handlers.NewErrorJson(ctx)
		}

		newFilename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
		savePath := fmt.Sprintf("./Asset/image/%s", newFilename)

		if err = ctx.SaveFile(file, savePath); err != nil {
			return handlers.NewErrorJson(ctx)
		}
		if currentImage != "" && currentImage != "default.jpg" {
			oldPath := fmt.Sprintf("./Asset/image/%s", currentImage)
			os.Remove(oldPath)
		}

		filename = newFilename
	}

	req := order_service.UpdateMenuReq{
		ID:         uint(id),
		CategoryID: uint(categoryID),
		Name:       name,
		Price:      price,
		Image:      filename,
		Status:     status,
	}

	err = o.orderService.UpdateMenu(req)
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessMsg(ctx)
}

func (o *orderController) DeleteMenuCtr(ctx *fiber.Ctx) error {
	req := new(order_service.DeleteMenuReq)
	err := ctx.BodyParser(req)
	if err != nil {
		return handlers.NewErrorJson(ctx)
	}
	errValid := validation.Validate(*req)
	if errValid != nil {
		return handlers.NewErrorValidate(ctx, errValid[0].Error)
	}
	err = o.orderService.DeleteMenu(*req)
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessMsg(ctx)
}
func (o *orderController) DeleteOrderItems(ctx *fiber.Ctx) error {
	req := new(order_service.DeleteOrderItemReq)
	err := ctx.BodyParser(req)
	if err != nil {
		return handlers.NewErrorJson(ctx)
	}
	errValid := validation.Validate(*req)
	if errValid != nil {
		return handlers.NewErrorValidate(ctx, errValid[0].Error)
	}
	err = o.orderService.DeleteOrderItem(*req)
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessMsg(ctx)
}

func (o *orderController) GetAllOrdersNoFilterCtr(ctx *fiber.Ctx) error {
	res, err := o.orderService.GetAllOrdersNoFilter()
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessResponse(ctx, res)
}

func NewOrderController(
	orderService *order_service.OrderService,
) OrderController {
	return &orderController{
		orderService: *orderService,
	}
}
