package order_service

type CreateCategoryReq struct {
	Name      string `json:"name" validate:"required"`
	CreatedAt string `json:"created_at"`
}

type UpdateCategoryReq struct {
	ID   uint   `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type DeleteCategoryReq struct {
	ID uint `json:"id" validate:"required"`
}
type DeleteOrderItemReq struct {
	CustomerID uint `json:"customer_id"`
	MenuID     uint `json:"menu_id"`
}

type CreateMenuReq struct {
	CategoryID uint   `json:"category_id" validate:"required"`
	Name       string `json:"name" `
	Image      string `json:"image" `
	Price      int    `json:"price" `
	Status     bool   `json:"status" `
}

type UpdateMenuReq struct {
	ID         uint   `json:"id" validate:"required"`
	CategoryID uint   `json:"category_id" `
	Name       string `json:"name"`
	Image      string `json:"image"`
	Price      int    `json:"price"`
	Status     bool   `json:"status"`
}

type DeleteMenuReq struct {
	ID uint `json:"id" validate:"required"`
}

type MenuResponse struct {
	ID           uint   `json:"id"`
	CategoryID   uint   `json:"category_id"`
	Name         string `json:"name"`
	Image        string `json:"image"`
	Price        int    `json:"price"`
	Status       bool   `json:"status"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	CategoryName string `json:"category_name"`
}

type UpdateOrderReq struct {
	ID          uint   `json:"id" validate:"required"`
	OrderStatus string `json:"order_status"`
	UpdatedAt   string `json:"updated_at"`
}

type DeleteOrderReq struct {
	ID uint `json:"id" validate:"required"`
}

type OrderReq struct {
	CustomerID uint `json:"customer_id" validate:"required"`
}
type OrderItemUpdateStatusReq struct {
	OrderItemID uint   `json:"order_item_id" validate:"required"`
	OrderStatus string `json:"order_status" validate:"required"`
}

//type OrderUpdateStatusReq struct {
//	ID         uint   `json:"id" validate:"required"`
//	MenuStatus string `json:"menu_status" `
//}

//	type CreateOrder struct {
//		CustomerID *uint  `json:"customer_id" validate:"required"`
//		MenuID     *uint  `json:"menu_id" validate:"required"`
//		Quantity   int    `json:"quantity" validate:"required"`
//		CreatedAt  string `json:"created_at"`
//		UpdatedAt  string `json:"updated_at"`
//		MenuStatus string `json:"menu_status"`
//	}
type CreateOrder struct {
	CustomerID *uint        `json:"customer_id" validate:"required"`
	Items      []CreateItem `json:"items" validate:"required,dive"`
}

type CreateItem struct {
	MenuID   *uint `json:"menu_id" validate:"required"`
	Quantity int   `json:"quantity" validate:"required"`
}

type OrderItemResponse struct {
	ID           uint   `json:"id"`
	CustomerID   uint   `json:"customer_id"`
	TableNumber  string `json:"table_number"`
	MenuID       uint   `json:"menu_id"`
	Name         string `json:"name"`
	Image        string `json:"image"`
	CategoryID   uint   `json:"category_id"`
	Price        int    `json:"price"`
	Quantity     int    `json:"quantity"`
	CategoryName string `json:"category_name"`
}

type OrderResponse struct {
	ID          uint                `json:"id"`
	CustomerID  uint                `json:"customer_id"`
	OrderStatus string              `json:"order_status"`
	TableNumber string              `json:"table_number"`
	CreatedAt   string              `json:"created_at"`
	OrderItems  []OrderItemResponse `json:"order_items"`
}

type CategoryResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
type CreateOrderRequest struct {
	CustomerID uint                     `json:"customer_id"`
	OrderItems []CreateOrderItemRequest `json:"order_items"`
}

type CreateOrderItemRequest struct {
	MenuID   uint `json:"menu_id"`
	Quantity int  `json:"quantity"`
}
