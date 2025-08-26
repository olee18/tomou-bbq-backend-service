package repositories

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type RedisRepository interface {
	GetDataRedis(key string) (string, error)
	SetDataRedis(key, value string) error
	DelDataRedis(key string) error

	GetHashRedis(keyStr, key string) (string, error)
	SetHasRedis(keyStr, key, value string) error
	DelHasRedis(keyStr, key string) error
}

type redisRepository struct{ rd *redis.Client }

func (r redisRepository) GetDataRedis(key string) (string, error) {
	result, err := r.rd.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (r redisRepository) SetDataRedis(key, value string) error {
	err := r.rd.Set(context.Background(), key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}
func (r redisRepository) DelDataRedis(key string) error {
	err := r.rd.Del(context.Background(), key).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r redisRepository) GetHashRedis(keyStr, key string) (string, error) {
	result, err := r.rd.HGet(context.Background(), keyStr, key).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (r redisRepository) SetHasRedis(keyStr, key, value string) error {
	err := r.rd.HSet(context.Background(), keyStr, key, value).Err()
	if err != nil {
		return err
	}
	return nil
}
func (r redisRepository) DelHasRedis(keyStr, key string) error {
	err := r.rd.HDel(context.Background(), keyStr, key).Err()
	if err != nil {
		return err
	}
	return nil
}

func NewRedisRepository(rd *redis.Client) RedisRepository {
	return &redisRepository{rd: rd}
}
