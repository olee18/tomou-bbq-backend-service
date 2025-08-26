package middlewares

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	jtoken "github.com/golang-jwt/jwt/v4"
	"laotop_final/config"
	"time"
)

func GenerateTokenWeb(id, roleId string, permission []string) (string, error) {
	claims := jtoken.MapClaims{
		"id":         id,
		"role_id":    roleId,
		"permission": permission,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(config.Env("jwt.web")))
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func GenerateTokenApi(id, phone string) (string, error) {
	claims := jtoken.MapClaims{
		"id":      id,
		"phone":   phone,
		"exp":     time.Now().Add(time.Hour * 720).Unix(),
		"Issuer":  id,
		"Subject": id,
	}
	token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(config.Env("jwt.api")))
	if err != nil {
		return "", err
	}
	return accessToken, nil
}
func AuthorizationCheckWeb(ctx *fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(config.Env("jwt.web")),
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			if err != nil {
				return err
			}
			return nil
		},
	})(ctx)
}

func AuthorizationCheckApi(ctx *fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(config.Env("jwt.api")),
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			if err != nil {
				return err
			}
			return nil
		},
	})(ctx)
}
func GetCustomerID(ctx *fiber.Ctx) string {
	user := ctx.Locals("user").(*jtoken.Token)
	claims := user.Claims.(jtoken.MapClaims)
	id, ok := claims["id"].(string)
	if !ok {
		return ""
	}
	return id
}
func GetCustomerPhone(ctx *fiber.Ctx) string {
	user := ctx.Locals("user").(*jtoken.Token)
	claims := user.Claims.(jtoken.MapClaims)
	id, ok := claims["phone"].(string)
	if !ok {
		return ""
	}
	return id
}
func GetUserID(ctx *fiber.Ctx) string {
	user := ctx.Locals("user").(*jtoken.Token)
	claims := user.Claims.(jtoken.MapClaims)
	id, ok := claims["id"].(string)
	if !ok {
		return ""
	}
	return id
}
func GetUserRoleID(ctx *fiber.Ctx) string {
	user := ctx.Locals("user").(*jtoken.Token)
	claims := user.Claims.(jtoken.MapClaims)
	id, ok := claims["role_id"].(string)
	if !ok {
		return ""
	}
	return id
}
func GetUserPermission(c *fiber.Ctx) []interface{} {
	user := c.Locals("user").(*jtoken.Token)
	claims := user.Claims.(jtoken.MapClaims)
	permissionClaim, ok := claims["permission"].([]interface{})
	if !ok {
		return []interface{}{}
	}
	return permissionClaim
}
