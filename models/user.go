package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Username     string    `json:"username" gorm:"unique"`
	AccessStatus *bool     `json:"access_status"`
	Password     string    `json:"password"`
	RoleID       int       `json:"role_id" `
	Role         Role      `json:"role" gorm:"foreignKey:RoleID"`
	LoginAt      time.Time `json:"login_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Role struct {
	ID         int          `json:"id"`
	Name       string       `json:"name" gorm:"unique"`
	Permission []Permission `json:"permission" gorm:"many2many:role_permissions"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}

type Permission struct {
	ID      int    `json:"id" gorm:"primaryKey"`
	Name    string `json:"name"  gorm:"unique"`
	Keyword string `json:"keyword"  gorm:"unique"`
	Sort    string `json:"sort"`
}

type RolePermission struct {
	RoleID       int `gorm:"primaryKey"`
	PermissionID int `gorm:"primaryKey"`
}
