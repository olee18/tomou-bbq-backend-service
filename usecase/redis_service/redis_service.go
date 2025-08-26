package redis_service

import (
	"laotop_final/repositories"
)

type RedisService interface {
	GetHashRedisWeb(keyStr, key string) (string, error)
	GetHashRedisApi(keyStr, key string) (string, error)
}

type redisService struct {
	repositoryRedis repositories.RedisRepository
}

func (s redisService) GetHashRedisWeb(keyStr, key string) (string, error) {
	result, err := s.repositoryRedis.GetHashRedis(keyStr, key)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (s redisService) GetHashRedisApi(keyStr, key string) (string, error) {
	result, err := s.repositoryRedis.GetHashRedis(keyStr, key)
	if err != nil {
		return "", err
	}
	return result, nil
}

func NewRedisService(
	repositoryRedis *repositories.RedisRepository,
	// repo
) RedisService {
	return &redisService{
		repositoryRedis: *repositoryRedis,
		//repo
	}
}
