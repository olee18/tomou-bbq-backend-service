package models

import "time"

type Order struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	CustomerID  uint        `json:"customer_id"`
	OrderStatus string      `json:"order_status"`
	TableNumber string      `json:"table_number"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	OrderItems  []OrderItem `gorm:"foreignKey:OrderID" json:"order_items"`
}

type OrderItem struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	OrderID    uint      `json:"order_id"`
	MenuID     uint      `json:"menu_id"`
	Name       string    `json:"name"`
	Image      string    `json:"image"`
	CategoryID uint      `json:"category_id"`
	Price      int       `json:"price"`
	Quantity   int       `json:"quantity"`
	Menu       Menu      `gorm:"foreignKey:MenuID" json:"menu"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Order      Order     `gorm:"foreignKey:OrderID" json:"order"`
	Category   *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

type Menu struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	CategoryID uint      `json:"category_id"`
	Name       string    `json:"name"`
	Image      string    `json:"image"`
	Price      int       `json:"price"`
	Status     bool      `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Category   Category  `gorm:"foreignKey:CategoryID"`
}

type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
