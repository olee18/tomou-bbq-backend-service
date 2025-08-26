package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"laotop_final/handlers"
)

func Permission(key string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		permission := GetUserPermission(c)
		included := includes(permission, key)
		if included {
			return c.Next()
		}
		return handlers.NewErrorPermission(c)
	}
}
func includes(permission []interface{}, key string) bool {
	for _, item := range permission {
		if item == key {
			return true
		}
	}
	return false
}
