package web_route

import (
	"laotop_final/handlers/web_handler"
	"laotop_final/middlewares"
	"laotop_final/routes"

	"github.com/gofiber/fiber/v2"
)

type webRoutes struct {
	userCtr     web_handler.UserController
	redisCtr    middlewares.RedisController
	customerCtr web_handler.CustomerController
	orderCtr    web_handler.OrderController
	billsCtr    web_handler.BillsController
}

func (w webRoutes) Install(app *fiber.App) {
	route := app.Group("web/", func(ctx *fiber.Ctx) error {
		return ctx.Next()
	})
	route.Post("sign-in", w.userCtr.SignIn)

	routes := app.Group("web/auth/", func(ctx *fiber.Ctx) error {
		return ctx.Next()
	})
	//user
	routes.Post("infos-get", w.userCtr.GetUserInfos)
	routes.Post("password-change", w.userCtr.ChangePassword)
	routes.Post("profile-update", w.userCtr.GetUserInfos)
	routes.Post("sign-out", w.userCtr.SignOut)

	//permission
	routes.Post("permission-get", w.userCtr.GetPermission)
	//role
	routes.Post("role-get", w.userCtr.GetRole)
	routes.Post("role-create", w.userCtr.CreatedRole)
	routes.Post("role-update", w.userCtr.UpdateRole)
	routes.Post("role-delete", w.userCtr.DeleteRole)

	routes.Post("user-get", w.userCtr.GetUser)
	routes.Post("user-create", w.userCtr.CreateUser)
	routes.Post("user-update", w.userCtr.UpdateUser)
	routes.Post("user-delete", w.userCtr.DeleteUser)

	routes.Post("customer-get", w.customerCtr.GetCustomerCtr)
	routes.Post("customer-create", w.customerCtr.CreateCustomerCtr)
	routes.Post("customer-update", w.customerCtr.UpdateCustomerCtr)
	routes.Post("customer-delete", w.customerCtr.DeleteCustomerCtr)
	routes.Post("customer-by-id", w.customerCtr.GetCustomerByIDCtr)

	routes.Post("order-get", w.orderCtr.GetOrderCtr)
	routes.Post("order-get-by-table-id", w.orderCtr.GetORderByCustomerID)
	routes.Post("order-create", w.orderCtr.CreateOrderCtr)
	routes.Post("order-update", w.orderCtr.UpdateOrderCtr)
	routes.Post("order-delete", w.orderCtr.DeleteOrderCtr)

	routes.Post("category-id", w.orderCtr.GetCategoryByIDCtr)
	routes.Post("category-get", w.orderCtr.GetCategoryCtr)
	routes.Post("category-create", w.orderCtr.CreateCategoryCtr)
	routes.Post("category-update", w.orderCtr.UpdateCategoryCtr)
	routes.Post("category-delete", w.orderCtr.DeleteCategoryCtr)

	routes.Post("menu-id", w.orderCtr.GetMenuByCategoryIDCtr)
	routes.Post("menu-get", w.orderCtr.GetMenuCtr)
	routes.Post("menu-create", w.orderCtr.CreateMenuCtr)
	routes.Post("menu-update", w.orderCtr.UpdateMenuCtr)
	routes.Post("menu-delete", w.orderCtr.DeleteMenuCtr)

	//bill history //
	route.Post("get-bills-by-id", w.billsCtr.GetBillsCtr)
	route.Post("get-bill-history", w.billsCtr.GetBillHistroy)
	route.Post("crete-bill-by-table", w.billsCtr.InsertBilssCtr)
	route.Post("update-bills", w.billsCtr.UpdateBillCtr)
	route.Post("bills-delete", w.billsCtr.DeleteBillCtr)
	route.Post("bills-customer-id", w.billsCtr.ShowBillClientByCustomerIDCtr)

	route.Post("report-bill", w.billsCtr.GetBillitem)
	route.Post("report-order", w.orderCtr.GetAllOrderItemsCtr)
	route.Post("report-order-all", w.orderCtr.GetAllOrdersNoFilterCtr)
	route.Post("delete-bill-item", w.billsCtr.DeleteBillItemCtr)
	route.Post("delete-order-item", w.orderCtr.DeleteOrderItems)
}

func NewWebRoutes(

	userCtr *web_handler.UserController,
	redisCtr *middlewares.RedisController,
	customerCtr web_handler.CustomerController,
	orderCtr *web_handler.OrderController,
	billsCtr *web_handler.BillsController,

) routes.Routes {
	return &webRoutes{
		userCtr:     *userCtr,
		redisCtr:    *redisCtr,
		customerCtr: customerCtr,
		orderCtr:    *orderCtr,
		billsCtr:    *billsCtr,
	}
}
