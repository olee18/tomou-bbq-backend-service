package middlewares

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"laotop_final/handlers"
	"laotop_final/logs"
	"laotop_final/usecase/redis_service"
	"strings"
)

type RedisController interface {
	// GetAuthorizationWeb Insert your function interface
	GetAuthorizationWeb(ctx *fiber.Ctx) error
	GetAuthorizationApi(ctx *fiber.Ctx) error
}

type redisController struct {
	serviceRedis redis_service.RedisService
}

func (c redisController) GetAuthorizationWeb(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return handlers.NewErrorAuthResponse(ctx) // Stop here with 401
	}
	// Validate token
	err := AuthorizationCheckWeb(ctx)
	if err != nil {
		logs.Error(err)
		return handlers.NewErrorAuthResponse(ctx) // Stop here with 401
	}
	newAccessToken := strings.TrimPrefix(authHeader, "Bearer ")
	// Token checks in Redis
	roleID := GetUserRoleID(ctx)
	userID := GetUserID(ctx)
	getAuth, _ := c.serviceRedis.GetHashRedisWeb(roleID, userID)
	fmt.Println(newAccessToken != getAuth)
	if newAccessToken != getAuth {
		return handlers.NewErrorAuthResponse(ctx) // Stop here with 401
	}

	return ctx.Next() // Allow further middleware or route handlers
}

func (c redisController) GetAuthorizationApi(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return handlers.NewErrorAuthResponse(ctx)
	}
	newAccessToken := strings.TrimPrefix(authHeader, "Bearer ")
	err := AuthorizationCheckApi(ctx)
	if err != nil {
		logs.Error(err)
		return handlers.NewErrorAuthResponse(ctx)
	}
	data := GetCustomerID(ctx)
	getAuth, _ := c.serviceRedis.GetHashRedisApi("customers", data)
	if newAccessToken != getAuth {
		logs.Error("error")
		return handlers.NewErrorAuthResponse(ctx)
	}
	return ctx.Next()
}

func NewRedisController(
	serviceRedis *redis_service.RedisService,
) RedisController {
	return &redisController{
		serviceRedis: *serviceRedis,
	}
}
