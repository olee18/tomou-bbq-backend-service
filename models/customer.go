package models

import "time"

type DataCustomer struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	TableNumber   string    `json:"table_number"`
	Adults        int       `json:"adults"`
	Children      int       `json:"children"`
	TotalCustomer int       `json:"total_customer"`
	CanOrder      bool      `gorm:"default:true" json:"can_order"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Status        bool      `json:"status"`
}
