package models

import "time"

type Bill struct {
	ID            uint         `json:"id" gorm:"primaryKey"`
	OrderID       uint         `json:"order_id"`
	Order         Order        `json:"order" gorm:"foreignKey:OrderID;references:ID"`
	CustomerID    uint         `json:"customer_id"`
	DataCustomer  DataCustomer `json:"data_customer" gorm:"foreignKey:CustomerID;references:ID"`
	Name          string       `json:"name"`
	TotalAmount   int          `json:"total_amount"`
	TableNumber   string       `json:"table_number"`
	Adults        int          `json:"adults"`
	Children      int          `json:"children"`
	TotalCustomer int          `json:"total_customer"`
	Date          time.Time    `json:"date"`
	BillItems     []BillItem   `json:"bills_items" gorm:"foreignKey:BillID"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
	StatusBills   bool         `json:"bill_status"`
}

type BillItem struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	BillID      uint      `json:"bill_id"`
	OrderID     uint      `json:"order_id"`
	CustomerID  uint      `json:"customer_id" gorm:"-"`
	TableNumber string    `json:"table_number" gorm:"-"`
	CategoryID  uint      `json:"category_id"`
	MenuName    string    `json:"name"`
	Quantity    int       `json:"quantity"`
	Price       int       `json:"price"`
	TotalPrice  int       `json:"total_price"`
	Order       Order     `gorm:"foreignKey:OrderID;references:ID" json:"order"`
	Category    *Category `gorm:"foreignKey:CategoryID;references:ID" json:"category,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
