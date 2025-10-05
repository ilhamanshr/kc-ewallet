package service

import "github.com/gomodule/redigo/redis"

type RedisService struct {
	Pool *redis.Pool
}

func NewRedisService(pool *redis.Pool) *RedisService {
	return &RedisService{
		Pool: pool,
	}
}
