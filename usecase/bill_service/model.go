package bill_service

import "time"

type BillItemResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Quantity   int    `json:"quantity"`
	Price      int    `json:"price"`
	TotalPrice int    `json:"total_price"`
}
type BillItemWithOrder struct {
	ID         uint
	BillID     uint
	OrderID    uint
	CategoryID uint
	MenuName   string
	Quantity   int
	Price      int
	TotalPrice int
	CreatedAt  time.Time
	UpdatedAt  time.Time

	CustomerID  uint
	TableNumber string
}

type BillResponse struct {
	ID                  int                `json:"id" gorm:"primaryKey"`
	CustomerID          uint               `json:"customer_id"`
	TableNumber         string             `json:"table_number"`
	OrderID             uint               `json:"order_id"`
	TotalAmount         int                `json:"total_amount"`
	TotalCustomer       int                `json:"total_customer"`
	Adults              int                `json:"adults"`
	Children            int                `json:"children"`
	Date                string             `json:"date"`
	CreatedAt           time.Time          `json:"created_at"`
	UpdatedAt           time.Time          `json:"updated_at"`
	Items               []BillItemResponse `json:"items"`
	StatusBills         bool               `json:"bill_status"`
	TotalCombinedAmount int                `json:"total_combined_amount"`
}

type UpdateBillReq struct {
	ID                  int                `json:"id" binding:"required"`
	TableNumber         string             `json:"table_number" `
	TotalAmount         int                `json:"total_amount" `
	TotalCustomer       int                `json:"total_customer" `
	Adults              int                `json:"adults" `
	Children            int                `json:"children" `
	Items               []BillItemResponse `json:"items"`
	StatusBills         bool               `json:"bill_status"`
	TotalCombinedAmount int                `json:"total_combined_amount"`
	CreatedAt           string             `json:"created_at"`
	UpdatedAt           string             `json:"updated_at"`
}

type BillsReq struct {
	CustomerID  uint   `json:"customer_id" `
	TableNumber string `json:"table_number"`
	ID          int    `json:"id" gorm:"primaryKey"`
}
type BillsIDReq struct {
	CustomerID  uint   `json:"customer_id" `
	TableNumber string `json:"table_number"`
	ID          int    `json:"id" gorm:"primaryKey"`
}

type BillDeleteReq struct {
	ID uint `json:"id" validate:"required"`
}
type BillitemDeleteReq struct {
	ID uint `json:"id" validate:"required"`
}

type OrderResponse struct {
	ID           uint      `json:"id"`
	TableNumber  string    `json:"table_number"`
	CustomerID   *uint     `json:"customer_id" `
	CategoryID   uint      `json:"category_id"`
	CategoryName string    `json:"category_name"`
	MenuID       uint      `json:"menu_id"`
	Name         string    `json:"name"`
	Image        string    `json:"image"`
	Price        int       `json:"price"`
	Quantity     int       `json:"quantity"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
}
