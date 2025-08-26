package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"laotop_final/config"
	"laotop_final/database"
	"laotop_final/handlers/web_handler"
	"laotop_final/logs"
	"laotop_final/middlewares"
	"laotop_final/repositories"
	web2 "laotop_final/routes/web_route"
	"laotop_final/usecase/bill_service"
	"laotop_final/usecase/customer_service"
	"laotop_final/usecase/order_service"
	"laotop_final/usecase/redis_service"
	"laotop_final/usecase/user_service"
	"log"
)

func main() {
	//----------------------------------database----------------------------------//
	fmt.Println("=============================")
	postgresConnection, err := database.PostgresConnection()
	if err != nil {
		logs.Error(err)
		return
	}
	defer database.CloseConnectionPostgres(postgresConnection)
	fmt.Println("=============================")
	redisConnection, err := database.RedisConnection()
	if err != nil {
		logs.Error(err)
		return
	}
	defer database.CloseConnectionRedis(redisConnection)
	fmt.Println("=============================")
	//----------------------------------database----------------------------------//

	//----------------------------------config----------------------------------//
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	app.Use(logger.New())
	app.Use(cors.New())
	//note for allow client read file in this project
	app.Static("/image", "./Asset/image")
	//----------------------------------config----------------------------------//

	//----------------------------------repository----------------------------------//
	redisRepo := repositories.NewRedisRepository(redisConnection)

	userRepo := repositories.NewUserRepository(postgresConnection)
	customerRepo := repositories.NewCustomerRepository(postgresConnection)
	orderRepo := repositories.NewOrderRepository(postgresConnection)
	billsRepo := repositories.NewBillRePositroy(postgresConnection)
	//----------------------------------repository----------------------------------//
	//----------------------------------partner----------------------------------//
	//----------------------------------partner----------------------------------//

	//----------------------------------service----------------------------------//
	redisSrv := redis_service.NewRedisService(&redisRepo)
	userSrv := user_service.NewUserService(&userRepo, &redisRepo)
	customerSrv := customer_service.NewCustomerService(&customerRepo)
	orderSrv := order_service.NewOrderService(postgresConnection, &orderRepo, &customerRepo)
	billSrv := bill_service.NewBillsService(billsRepo, &customerRepo, &orderRepo, postgresConnection)
	//----------------------------------service----------------------------------//

	//----------------------------------middleware----------------------------------//
	redisMdw := middlewares.NewRedisController(&redisSrv)
	//----------------------------------middleware----------------------------------//
	//----------------------------------handler api----------------------------------//

	//----------------------------------handler api----------------------------------//

	//----------------------------------handler web----------------------------------//
	userWeb := web_handler.NewUserController(&userSrv)
	customerWeb := web_handler.NewCustomerController(&customerSrv)
	orderWeb := web_handler.NewOrderController(&orderSrv)
	billweb := web_handler.NewBillsController(&billSrv)
	//----------------------------------handler web----------------------------------//

	//----------------------------------routes----------------------------------//
	newWebRoute := web2.NewWebRoutes(
		&userWeb,
		&redisMdw,
		customerWeb,
		&orderWeb,
		&billweb,
	)

	//----------------------------------routes----------------------------------//
	newWebRoute.Install(app)
	log.Fatal(app.Listen(fmt.Sprintf(":%s", config.Env("app.port"))))

}
