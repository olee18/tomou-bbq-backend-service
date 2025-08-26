package database

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"laotop_final/config"
	"laotop_final/logs"
)

func RedisConnection() (*redis.Client, error) {
	redisHost := config.Env("redis.host")
	redisPORT := config.Env("redis.port")
	fmt.Println("Redis connecting")
	client := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPORT,
		Password: "",
		DB:       0,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	fmt.Println("Redis connected")
	return client, nil
}

func CloseConnectionRedis(client *redis.Client) {
	if client != nil {
		err := client.Close()
		if err != nil {
			logs.Error(err)
			return
		}
	}
}
